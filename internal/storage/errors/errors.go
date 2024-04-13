package errors

import "errors"

var (
	ErrExpressionNotFound   = errors.New("expression not found")
	ErrExpressionInProgress = errors.New("expression in progress")
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)
