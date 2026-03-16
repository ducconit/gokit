package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"

	"gopkg.in/gomail.v2"
)

type smtpMailer struct {
	cfg SMTPConfig
}

// NewSMTP creates a new SMTP mailer using gomail.v2.
func NewSMTP(cfg SMTPConfig) Mailer {
	return &smtpMailer{cfg: cfg}
}

func (s *smtpMailer) Send(ctx context.Context, msg Message) error {
	if s.cfg.Host == "" || s.cfg.Port == 0 {
		return ErrInvalidConfig
	}

	m := gomail.NewMessage()

	from := msg.From
	if from == "" {
		from = s.cfg.From
	}
	m.SetHeader("From", from)
	m.SetHeader("To", msg.To...)
	if len(msg.Cc) > 0 {
		m.SetHeader("Cc", msg.Cc...)
	}
	if len(msg.Bcc) > 0 {
		m.SetHeader("Bcc", msg.Bcc...)
	}
	m.SetHeader("Subject", msg.Subject)

	contentType := msg.ContentType
	if contentType == "" {
		contentType = "text/plain"
	}
	m.SetBody(contentType, msg.Body)

	for _, att := range msg.Attachments {
		m.Attach(att.Name, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, att.Data)
			return err
		}), gomail.SetHeader(map[string][]string{
			"Content-Type": {att.ContentType},
		}))
	}

	d := gomail.NewDialer(s.cfg.Host, s.cfg.Port, s.cfg.Username, s.cfg.Password)
	if s.cfg.TLS {
		d.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         s.cfg.Host,
		}
	}

	// gomail doesn't support context natively for DialAndSend, but we can check ctx before sending
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	return nil
}
