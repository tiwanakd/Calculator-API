package models

import "errors"

var (
	ErrNotFound           = errors.New("models: no matching record found")
	ErrInvalidCredentails = errors.New("models: invalid credentails")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)
