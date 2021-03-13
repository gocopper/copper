package cauth

import (
	"context"

	"github.com/tusharsoni/copper/cerrors"
	"gorm.io/gorm"
)

// ErrNotFound is returned when a model does not exist in the repository
var ErrNotFound = gorm.ErrRecordNotFound

// NewRepo instantiates and returns a new Repo.
func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// Repo represents the repository layer for all models in the cauth package.
type Repo struct {
	db *gorm.DB
}

// GetUserByUsername queries the users table for a user with the given username.
func (r *Repo) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User

	err := r.db.WithContext(ctx).Where("username=?", username).First(&user).Error
	if err != nil {
		return nil, cerrors.New(err, "failed to query user", map[string]interface{}{
			"username": username,
		})
	}

	return &user, nil
}

// SaveUser saves or creates the given user.
func (r *Repo) SaveUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// GetSession queries the sessions table for a session with the given uuid.
func (r *Repo) GetSession(ctx context.Context, uuid string) (*Session, error) {
	var session Session

	err := r.db.WithContext(ctx).Where("uuid=?", uuid).First(&session).Error
	if err != nil {
		return nil, cerrors.New(err, "failed to query session", map[string]interface{}{
			"uuid": uuid,
		})
	}

	return &session, nil
}

// SaveSession saves or creates the given session.
func (r *Repo) SaveSession(ctx context.Context, session *Session) error {
	return r.db.WithContext(ctx).Save(session).Error
}
