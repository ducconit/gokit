package logging

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDisabled(t *testing.T) {
	_, closer, err := New(Config{Mode: ModeDisabled})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := closer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestNewFile(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "app.log")

	logger, closer, err := New(Config{
		Mode:     ModeFile,
		Level:    "info",
		FilePath: fp,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	logger.Info().Msg("hello")
	if err := closer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if _, err := os.Stat(fp); err != nil {
		t.Fatalf("Stat: %v", err)
	}
}
