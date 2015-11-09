// Package slog provides structured logging.
//
// There are many good logging packages available, and it is worth asking
// why the world needs another one.
//
// Here are some differentiators for this package. Not all of them are
// unique, but this is the only package (to date) that has all of them.
//
// 1. Log messages are not formatted using a printf style interface. Each
// log message should have a constant message, which makes it easier to
// filter and search for messages. Any variable information is passed as
// properties in the message (see the WithValue function).
//  doSometingWith(a, b)
//  log.Debug(ctx, "did something",
//      log.WithValue("a", a),
//      log.WithValue("b", b))
//
// 2. Uses an api that allows for multiple options and parameters to be
// logged in a single call. (See "Functional options for friendly APIs"
// by Dave Cheney http://goo.gl/l2KaW3).
//  if err := doSometing(ctx, a); err != nil {
//      return log.Error(ctx, "cannot do someting",
//          log.WithValue("a", a),
//          log.WithError(err),
//          log.WithStatus(http.StatusBadRequest))
//  }
//
// 3. When a message is logged, a non-nil *Message value is returned, which
// can be returned as an error value.
//  if err := doSometing(); err != nil {
//      return log.Error(ctx, "cannot doSomething", log.WithError(err))
//  }
//
// 4. This package is context aware (golang.org/x/net/context). Contexts
// can be created with information that will be logged with the message.
//  ctx = log.NewContext(ctx, log.Property{"a", "SomeValue"})
//
//  // ... do some work and then
//
//  // the following message will include "a=SomeValue" from the context
//  log.Info(ctx, "some message")
//
// 5. By default messages are logged to stdout in logfmt format.
// (https://brandur.org/logfmt)
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
// Unlike Info, Warn and Error, this function does not return a pointer to the message.
// The reason is that Debug should never be used to return an error result.
func Debug(ctx context.Context, text string, opts ...Option) {
	Default.Debug(ctx, text, opts...)
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

// Set the output writer for the default logger.
func SetOutput(w io.Writer) {
	Default.SetOutput(w)
}

// Set the minimum log level for the default logger. By default
// the minimum log level is LevelInfo.
func SetMinLevel(level Level) {
	Default.SetMinLevel(level)
}

// Appends the handler to the list of handlers for the default logger.
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
