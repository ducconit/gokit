package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type CloserFunc func() error

func (f CloserFunc) Close() error { return f() }

func New(cfg Config) (zerolog.Logger, io.Closer, error) {
	level := zerolog.InfoLevel
	if cfg.Level != "" {
		lvl, err := zerolog.ParseLevel(cfg.Level)
		if err != nil {
			return zerolog.Logger{}, nil, err
		}
		level = lvl
	}

	switch cfg.Mode {
	case ModeDisabled:
		return zerolog.New(io.Discard).Level(zerolog.Disabled), CloserFunc(func() error { return nil }), nil
	case ModeConsole, ModeFile, ModeBoth, "":
	default:
		return zerolog.Logger{}, nil, fmt.Errorf("logging: invalid mode: %s", cfg.Mode)
	}

	var (
		consoleOut io.Writer
		fileOut    io.Writer
		closers    []io.Closer
	)

	if cfg.Mode == ModeConsole || cfg.Mode == ModeBoth || cfg.Mode == "" {
		if !cfg.DisableConsole {
			consoleOut = os.Stdout
			if cfg.ConsolePretty {
				consoleOut = zerolog.ConsoleWriter{Out: consoleOut}
			}
		}
	}

	if cfg.Mode == ModeFile || cfg.Mode == ModeBoth {
		if cfg.FilePath == "" {
			return zerolog.Logger{}, nil, fmt.Errorf("logging: missing file path")
		}
		f, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return zerolog.Logger{}, nil, err
		}
		fileOut = f
		closers = append(closers, f)
	}

	var out io.Writer
	switch {
	case consoleOut != nil && fileOut != nil:
		out = io.MultiWriter(consoleOut, fileOut)
	case consoleOut != nil:
		out = consoleOut
	case fileOut != nil:
		out = fileOut
	default:
		out = io.Discard
	}

	logger := zerolog.New(out).With().Timestamp().Logger().Level(level)

	return logger, CloserFunc(func() error {
		var firstErr error
		for i := len(closers) - 1; i >= 0; i-- {
			if err := closers[i].Close(); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		return firstErr
	}), nil
}
