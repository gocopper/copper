package cauth2

import "errors"

// ErrUserAlreadyExists is returned by UsersSvc when a user already exists. For example, signing up with an email
// with which a user already exists.
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrInvalidCredentials is returned by UsersSvc when the given credentials such as email/password or session token
// are incorrect.
var ErrInvalidCredentials = errors.New("invalid credentials")
