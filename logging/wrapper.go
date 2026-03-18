package logging

import (
	"io"

	"github.com/rs/zerolog"
)

// Logger is a generic logging interface.
type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, err error, fields ...any)
	With(fields ...any) Logger
}

// NewWrapper creates a new Logger with the given config.
func NewWrapper(cfg Config) (Logger, io.Closer, error) {
	logger, closer, err := New(cfg)
	if err != nil {
		return nil, nil, err
	}
	return &zerologWrapper{logger: logger}, closer, nil
}

// zerologWrapper implements Logger interface using zerolog.
type zerologWrapper struct {
	logger zerolog.Logger
}

func (w *zerologWrapper) Debug(msg string, fields ...any) {
	w.logger.Debug().Fields(fieldsToMap(fields)).Msg(msg)
}

func (w *zerologWrapper) Info(msg string, fields ...any) {
	w.logger.Info().Fields(fieldsToMap(fields)).Msg(msg)
}

func (w *zerologWrapper) Warn(msg string, fields ...any) {
	w.logger.Warn().Fields(fieldsToMap(fields)).Msg(msg)
}

func (w *zerologWrapper) Error(msg string, err error, fields ...any) {
	e := w.logger.Error()
	if err != nil {
		e = e.Err(err)
	}
	e.Fields(fieldsToMap(fields)).Msg(msg)
}

func (w *zerologWrapper) With(fields ...any) Logger {
	return &zerologWrapper{
		logger: w.logger.With().Fields(fieldsToMap(fields)).Logger(),
	}
}

func fieldsToMap(fields []any) map[string]any {
	if len(fields) == 0 {
		return nil
	}
	m := make(map[string]any)
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			if key, ok := fields[i].(string); ok {
				m[key] = fields[i+1]
			}
		}
	}
	return m
}
