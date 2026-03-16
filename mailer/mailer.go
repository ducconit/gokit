package mailer

import (
	"context"
	"io"
)

// Message represents an email message.
type Message struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	ContentType string // e.g. "text/plain" or "text/html"
	Attachments []Attachment
}

// Attachment represents an email attachment.
type Attachment struct {
	Name        string
	ContentType string
	Data        io.Reader
}

// Mailer defines the interface for sending emails.
type Mailer interface {
	Send(ctx context.Context, msg Message) error
}

// New creates a new Mailer from configuration.
func New(cfg Config) (Mailer, error) {
	if cfg.SMTP.Host != "" {
		return NewSMTP(cfg.SMTP), nil
	}
	return nil, ErrInvalidConfig
}
