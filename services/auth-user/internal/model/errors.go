package model

import "errors"

// Domain errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidRole       = errors.New("invalid role")
	ErrUserBlocked       = errors.New("user is blocked")

	ErrTokenNotFound = errors.New("refresh token not found")
	ErrTokenExpired  = errors.New("refresh token expired")
	ErrTokenRevoked  = errors.New("refresh token revoked")
	ErrInvalidToken  = errors.New("invalid refresh token")
)
