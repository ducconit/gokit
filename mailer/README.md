# Mailer Package

The `mailer` package provides a simple interface for sending emails. It currently supports SMTP.

## Usage

```go
package main

import (
	"context"
	"github.com/ducconit/gokit/mailer"
)

func main() {
	cfg := mailer.Config{
		SMTP: mailer.SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "user",
			Password: "password",
			From:     "noreply@example.com",
		},
	}

	m, err := mailer.New(cfg)
	if err != nil {
		panic(err)
	}

	msg := mailer.Message{
		To:      []string{"user@example.com"},
		Subject: "Hello",
		Body:    "Hello, world!",
	}

	err = m.Send(context.Background(), msg)
	if err != nil {
		panic(err)
	}
}
```

## Features

- SMTP support (PlainAuth, TLS)
- Multipart messages (HTML/Text body)
- Attachments
- Custom headers (Cc, Bcc, Subject)
