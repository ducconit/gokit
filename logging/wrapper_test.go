package logging

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewWrapperDisabled(t *testing.T) {
	_, closer, err := NewWrapper(Config{Mode: ModeDisabled})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := closer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestNewWrapperFile(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "app.log")

	logger, closer, err := NewWrapper(Config{
		Mode:     ModeFile,
		Level:    "info",
		FilePath: fp,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	logger.Info("hello")
	if err := closer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if _, err := os.Stat(fp); err != nil {
		t.Fatalf("Stat: %v", err)
	}
}
