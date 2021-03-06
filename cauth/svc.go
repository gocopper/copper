package cauth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tusharsoni/copper/cerrors"
	"github.com/tusharsoni/copper/crandom"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ErrInvalidCredentials is returned when a credential check fails. This usually happens during the login process.
var ErrInvalidCredentials = errors.New("invalid credentials")

// NewSvc instantiates and returns a new Svc.
func NewSvc(repo *Repo) *Svc {
	return &Svc{
		repo: repo,
	}
}

// Svc provides methods to manage users and sessions.
type Svc struct {
	repo *Repo
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
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// LoginParams hold the params needed to login a user.
type LoginParams struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// Signup creates a new user. If contact methods such as email or phone are provided, it will send verification
// codes so them. It creates a new session for this newly created user and returns that.
func (s *Svc) Signup(ctx context.Context, p SignupParams) (*SessionResult, error) {
	if p.Username == nil {
		return nil, cerrors.New(nil, "at least one login method is required", nil)
	}

	if p.Username != nil && p.Password == nil {
		return nil, cerrors.New(nil, "password is required with username", nil)
	}

	var hashedPassword []byte

	if p.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*p.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, cerrors.New(err, "failed to hash password", nil)
		}

		hashedPassword = hash
	}

	user := &User{
		UUID:               uuid.New().String(),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Username:           p.Username,
		Password:           hashedPassword,
		PasswordResetToken: nil,
	}

	err := s.repo.SaveUser(ctx, user)
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
	if p.Username == nil {
		return nil, cerrors.New(nil, "at least one login method is required", nil)
	}

	if p.Username != nil && p.Password == nil {
		return nil, cerrors.New(nil, "password is required with username", nil)
	}

	return s.loginWithUsernamePassword(ctx, *p.Username, *p.Password)
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
	if err != nil {
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
