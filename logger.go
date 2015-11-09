package slog

import (
	"io"
	"os"
	"sync"

	"golang.org/x/net/context"
)

const (
	// Bits or'ed together to control what is printed.
	LTimestamp = 1 << iota // Include timestamp in output
	LUTC                   // If LTimestamp is set, use UTC

	LLF // On Windows, print LF only, not CRLF. Required for testing

	LDefault = LTimestamp
)

type Logger interface {
	Debug(ctx context.Context, text string, opts ...Option) *Message
	Info(ctx context.Context, text string, opts ...Option) *Message
	Warn(ctx context.Context, text string, opts ...Option) *Message
	Error(ctx context.Context, text string, opts ...Option) *Message

	Flags() int
	SetFlags(flags int)

	NewWriter(ctx context.Context) io.Writer
	SetOutput(w io.Writer)
	SetMinLevel(level Level)
	AddHandler(h Handler)
}

type Handler interface {
	Handle(msgs []*Message)
}

type loggerImpl struct {
	mu       sync.Mutex // ensures atomic writes; protects the following fields
	out      io.Writer  // destination for output
	handlers []Handler  // list of handlers
	minLevel Level      // minimum level to log
	flags    int        // flag options
}

// New returns a new Logger with default settings. Writes to stdout, and
// does not print debug messages.
func New() Logger {
	return &loggerImpl{
		// NOTE: differs from std logging in that default is standard output
		// not standard error. This is consistent with 12 factor app, but
		// is it the appropriate default.
		out:      os.Stdout,
		minLevel: LevelInfo,
	}
}

func (l *loggerImpl) Debug(ctx context.Context, text string, opts ...Option) *Message {
	m := newMessage(ctx, LevelDebug, text)
	m.applyOpts(opts)
	l.output(m)
	return m
}

func (l *loggerImpl) Info(ctx context.Context, text string, opts ...Option) *Message {
	m := newMessage(ctx, LevelInfo, text)
	m.applyOpts(opts)
	l.output(m)
	return m
}

func (l *loggerImpl) Warn(ctx context.Context, text string, opts ...Option) *Message {
	m := newMessage(ctx, LevelWarning, text)
	m.applyOpts(opts)
	l.output(m)
	return m
}

func (l *loggerImpl) Error(ctx context.Context, text string, opts ...Option) *Message {
	m := newMessage(ctx, LevelError, text)
	m.applyOpts(opts)
	l.output(m)
	return m
}

func (l *loggerImpl) Flags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flags
}

func (l *loggerImpl) SetFlags(flags int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flags = flags
}

func (l *loggerImpl) NewWriter(ctx context.Context) io.Writer {
	return &writer{
		ctx:    ctx,
		logger: l,
	}
}

func (l *loggerImpl) SetMinLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

func (l *loggerImpl) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *loggerImpl) AddHandler(h Handler) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.handlers = append(l.handlers, h)
}

// output provides the common functionality to output a message.
func (l *loggerImpl) output(m *Message) {
	// TODO: if out is a tty, use ansi sequences to print color-coded output.
	// For now, just always write to the output using logfmt format.
	l.mu.Lock()
	defer l.mu.Unlock()

	if m.Level >= l.minLevel {
		if l.output != nil {
			buf := m.logfmtBuffer(l.flags)
			// <HACK>
			// We need this hack to print LF only so that the Example
			// code will work on Windows. Of CRLF is printed on Windows,
			// the example tests fail because the CR is not trimmed.
			if l.flags&LLF != 0 {
				buf.WriteNewLine()
			} else {
				buf.WriteEOL()
			}
			// </HACK>
			buf.WriteTo(l.out)
			buf.Reset()
		}

		// TODO: could reduce locking here by having a goroutine and a buffered
		// channel for each handler. The goroutine could read from the buffered
		// channel and group messages into a slice and send the slice to the
		// handler. This is why the Handler type accepts a slice of messages.
		// For now, the implementation is simple.
		messages := []*Message{m}
		for _, handler := range l.handlers {
			handler.Handle(messages)
		}
	}
}
