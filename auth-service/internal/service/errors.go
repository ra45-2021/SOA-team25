package service

import "errors"

var (
	ErrBadRole        = errors.New("role must be GUIDE or TOURIST")
	ErrHashFailed     = errors.New("hash failed")
	ErrIDGenFailed    = errors.New("id generation failed")
	ErrDuplicateUser  = errors.New("username or email already exists")
	ErrInvalidCreds   = errors.New("invalid credentials")
	ErrTokenFailed    = errors.New("token failed")
)
