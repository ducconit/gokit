package mailer

import (
	"context"
	"testing"
)

func TestMailerConfig(t *testing.T) {
	cfg := Config{
		SMTP: SMTPConfig{
			Host: "localhost",
			Port: 25,
		},
	}

	m, err := New(cfg)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if m == nil {
		t.Fatal("Expected mailer, got nil")
	}
}

func TestMailerInvalidConfig(t *testing.T) {
	cfg := Config{}
	_, err := New(cfg)
	if err == nil {
		t.Fatal("Expected error for invalid config, got nil")
	}
}

func TestSendWithContextCancel(t *testing.T) {
	cfg := Config{
		SMTP: SMTPConfig{
			Host: "localhost",
			Port: 25,
		},
	}
	m, _ := New(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := m.Send(ctx, Message{
		To:      []string{"to@example.com"},
		Subject: "Test",
		Body:    "Test",
	})

	if err == nil {
		t.Fatal("Expected error for cancelled context, got nil")
	}
}
