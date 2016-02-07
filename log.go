package slog

import (
	"io"

	"golang.org/x/net/context"
)

var (
	// Default is the default logger.
	Default = New()
)

// Debug logs a debug level message to the default logger.
// Returns a non-nil *Message, which can be used as an error value.
func Debug(ctx context.Context, text string, opts ...Option) *Message {
	return Default.Debug(ctx, text, opts...)
}

// Info logs an informational level message to the default logger.
// Returns a non-nil *Message, which can be used as an error value.
func Info(ctx context.Context, text string, opts ...Option) *Message {
	return Default.Info(ctx, text, opts...)
}

// Warn logs a warning level message to the default logger.
// Returns a non-nil *Message, which can be used as an error value.
func Warn(ctx context.Context, text string, opts ...Option) *Message {
	return Default.Warn(ctx, text, opts...)
}

// Error logs an error level message to the default logger.
// Returns a non-nil *Message, which can be used as an error value.
func Error(ctx context.Context, text string, opts ...Option) *Message {
	return Default.Error(ctx, text, opts...)
}

// SetOutput sets the output writer for the default logger.
func SetOutput(w io.Writer) {
	Default.SetOutput(w)
}

// SetMinLevel sets the minimum log level for the default logger. By default
// the minimum log level is LevelInfo.
func SetMinLevel(level Level) {
	Default.SetMinLevel(level)
}

// AddHandler appends the handler to the list of handlers for the default logger.
func AddHandler(h Handler) {
	Default.AddHandler(h)
}

// NewWriter creates a new writer that can be used to integrate with the
// standard log package. The main use case for this is to log messages
// generated from the standard library, in particular the net/http package.
// See the example for more information.
func NewWriter(ctx context.Context) io.Writer {
	return Default.NewWriter(ctx)
}
