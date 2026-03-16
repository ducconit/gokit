package mailer

import "errors"

var (
	ErrInvalidConfig = errors.New("mailer: invalid configuration")
	ErrSendFailed    = errors.New("mailer: failed to send email")
)
