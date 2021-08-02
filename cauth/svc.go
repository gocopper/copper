package cauth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gocopper/copper/cconfig"

	"github.com/gocopper/copper/cmailer"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/crandom"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ErrInvalidCredentials is returned when a credential check fails. This usually happens during the login process.
var ErrInvalidCredentials = errors.New("invalid credentials")

// NewSvc instantiates and returns a new Svc.
func NewSvc(repo *Repo, mailer cmailer.Mailer, appConfig cconfig.Config) (*Svc, error) {
	var config Config

	err := appConfig.Load("cauth", &config)
	if err != nil {
		return nil, cerrors.New(err, "failed to load cauth config", nil)
	}

	return &Svc{
		repo:   repo,
		mailer: mailer,
		config: &config,
	}, nil
}

// Svc provides methods to manage users and sessions.
type Svc struct {
	repo   *Repo
	mailer cmailer.Mailer
	config *Config
}

// SessionResult is usually used when a new session is created. It holds the plain session token that can be used
// to authenticate the session as well as other related entities such as the user and session.
type SessionResult struct {
	User              *User    `json:"user"`
	Session           *Session `json:"session"`
	PlainSessionToken string   `json:"plain_session_token"`
}

// SignupParams hold the params needed to signup a new user.
type SignupParams struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// LoginParams hold the params needed to login a user.
type LoginParams struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// Signup creates a new user. If contact methods such as email or phone are provided, it will send verification
// codes so them. It creates a new session for this newly created user and returns that.
func (s *Svc) Signup(ctx context.Context, p SignupParams) (*SessionResult, error) {
	if p.Username != nil && p.Password != nil {
		return s.signupWithUsernamePassword(ctx, *p.Username, *p.Password)
	}

	if p.Email != nil {
		return s.signupWithEmailOTP(ctx, *p.Email)
	}

	return nil, errors.New("invalid signup params")
}

func (s *Svc) signupWithEmailOTP(ctx context.Context, email string) (*SessionResult, error) {
	verificationCode := strconv.Itoa(int(crandom.GenerateRandomNumericalCode(s.config.VerificationCodeLen)))

	hashedVerificationCode, err := bcrypt.GenerateFromPassword([]byte(verificationCode), bcrypt.DefaultCost)
	if err != nil {
		return nil, cerrors.New(err, "failed to hash verification code", nil)
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, cerrors.New(err, "failed to get user by email", map[string]interface{}{
			"email": email,
		})
	}

	if user == nil {
		user = &User{
			UUID:      uuid.New().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Email:     &email,
		}
	}

	user.Password = hashedVerificationCode

	err = s.repo.SaveUser(ctx, user)
	if err != nil {
		return nil, cerrors.New(err, "failed to save user", nil)
	}

	emailBody := fmt.Sprintf("Your verification code is %s", verificationCode)

	err = s.mailer.Send(ctx, cmailer.SendParams{
		From:      s.config.VerificationEmailFrom,
		To:        []string{email},
		Subject:   s.config.VerificationEmailSubject,
		PlainBody: &emailBody,
	})
	if err != nil {
		return nil, cerrors.New(err, "failed to send verification code email", map[string]interface{}{
			"to": email,
		})
	}

	return &SessionResult{
		User: user,
	}, nil
}

func (s *Svc) signupWithUsernamePassword(ctx context.Context, username, password string) (*SessionResult, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, cerrors.New(err, "failed to hash password", nil)
	}

	user := &User{
		UUID:      uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Username:  &username,
		Password:  hashedPassword,
	}

	err = s.repo.SaveUser(ctx, user)
	if err != nil {
		return nil, cerrors.New(err, "failed to save user", nil)
	}

	session, plainSessionToken, err := s.createSession(ctx, user.UUID)
	if err != nil {
		return nil, cerrors.New(err, "failed to create a session", map[string]interface{}{
			"userUUID": user.UUID,
		})
	}

	return &SessionResult{
		User:              user,
		Session:           session,
		PlainSessionToken: plainSessionToken,
	}, nil
}

// Login logs in an existing user with the given credentials. If the login succeeds, it creates a new session
// and returns it.
func (s *Svc) Login(ctx context.Context, p LoginParams) (*SessionResult, error) {
	if p.Username != nil && p.Password != nil {
		return s.loginWithUsernamePassword(ctx, *p.Username, *p.Password)
	}

	if p.Email != nil && p.Password != nil {
		return s.loginWithEmailPassword(ctx, *p.Email, *p.Password)
	}

	return nil, cerrors.New(nil, "invalid login params", nil)
}

func (s *Svc) loginWithEmailPassword(ctx context.Context, email, password string) (*SessionResult, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidCredentials
	} else if err != nil {
		return nil, cerrors.New(err, "failed to get user by email", map[string]interface{}{
			"email": email,
		})
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	session, plainSessionToken, err := s.createSession(ctx, user.UUID)
	if err != nil {
		return nil, cerrors.New(err, "failed to create session", map[string]interface{}{
			"userUUID": user.UUID,
		})
	}

	return &SessionResult{
		User:              user,
		Session:           session,
		PlainSessionToken: plainSessionToken,
	}, nil
}

func (s *Svc) loginWithUsernamePassword(ctx context.Context, username, password string) (*SessionResult, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidCredentials
	} else if err != nil {
		return nil, cerrors.New(err, "failed to get user by username", map[string]interface{}{
			"username": username,
		})
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	session, plainSessionToken, err := s.createSession(ctx, user.UUID)
	if err != nil {
		return nil, cerrors.New(err, "failed to create session", map[string]interface{}{
			"userUUID": user.UUID,
		})
	}

	return &SessionResult{
		User:              user,
		Session:           session,
		PlainSessionToken: plainSessionToken,
	}, nil
}

func (s *Svc) createSession(ctx context.Context, userUUID string) (*Session, string, error) {
	const tokenLen = 128

	plainToken := crandom.GenerateRandomString(tokenLen)

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(plainToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", cerrors.New(err, "failed to hash session token", nil)
	}

	session := &Session{
		UUID:      uuid.New().String(),
		CreatedAt: time.Now(),
		UserUUID:  userUUID,
		Token:     hashedToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	err = s.repo.SaveSession(ctx, session)
	if err != nil {
		return nil, "", cerrors.New(err, "failed to create a new session", nil)
	}

	return session, plainToken, nil
}

// ValidateSession validates whether the provided plainToken is valid for the session identified by the given
// sessionUUID.
func (s *Svc) ValidateSession(ctx context.Context, sessionUUID, plainToken string) (bool, *Session, error) {
	session, err := s.repo.GetSession(ctx, sessionUUID)
	if err != nil && errors.Is(err, ErrNotFound) {
		return false, nil, nil
	} else if err != nil {
		return false, nil, cerrors.New(err, "failed to get session", map[string]interface{}{
			"sessionUUID": sessionUUID,
		})
	}

	err = bcrypt.CompareHashAndPassword(session.Token, []byte(plainToken))
	if err != nil {
		return false, nil, nil
	}

	return true, session, nil
}

// Logout invalidates the session identified by the given sessionUUID.
func (s *Svc) Logout(ctx context.Context, sessionUUID string) error {
	session, err := s.repo.GetSession(ctx, sessionUUID)
	if err != nil {
		return cerrors.New(err, "failed to get session", map[string]interface{}{
			"sessionUUID": sessionUUID,
		})
	}

	session.ExpiresAt = time.Now()

	err = s.repo.SaveSession(ctx, session)
	if err != nil {
		return cerrors.New(err, "failed to save session", map[string]interface{}{
			"sessionUUID": sessionUUID,
		})
	}

	return nil
}
