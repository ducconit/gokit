package auth

import "fmt"

var (
	ErrUnauthorized    = fmt.Errorf("unauthorized")
	ErrForbidden       = fmt.Errorf("forbidden")
	ErrSessionNotFound = fmt.Errorf("session not found")
)
