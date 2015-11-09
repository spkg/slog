package slog

import (
	"time"

	"golang.org/x/net/context"

	"github.com/spkg/slog/logfmt"
)

// Message contains all of the log message information.
// Note that *Message implements the error interface.
type Message struct {
	Timestamp  time.Time
	Level      Level
	Text       string
	Err        error
	Properties []Property
	Context    []Property
	code       string
	status     int
}

// Property is a key value pair associated with a Message.
type Property struct {
	Key   string
	Value interface{}
}

func newMessage(ctx context.Context, level Level, text string) *Message {
	m := &Message{
		Timestamp: time.Now(),
		Level:     level,
		Text:      text,
	}
	for data := fromContext(ctx); data != nil; data = data.Prev {
		m.Context = append(m.Context, Property{data.Key, data.Value})
	}

	return m
}

func (m *Message) applyOpt(opt Option) *Message {
	opt(m)
	return m
}

// applyOpts applies all of the option functions to the message.
func (m *Message) applyOpts(opts []Option) *Message {
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Error implements the error interface
func (m *Message) Error() string {
	return m.Text
}

// Code returns the code associated with the message.
// Implements the Coder interface.
func (m *Message) Code() string {
	return m.code
}

// SetCode sets the Code property.
func (m *Message) SetCode(code string) {
	m.code = code
}

// Status returns any status code associated with the message.
// This is intended to be a HTTP status code, but the application can
// use any numbering scheme.
func (m *Message) Status() int {
	return m.status
}

// SetStatus sets the status code associated with the message.
func (m *Message) SetStatus(status int) {
	m.status = status
}

// Logfmt writes the contents of the message to the buffer in logfmt format.
// See https://brandur.org/logfmt for a description of logfmt. Returns number
// of bytes written to w.
func (m *Message) Logfmt() string {
	buf := m.logfmtBuffer()
	s := buf.String()
	buf.Reset()
	return s
}

func (m *Message) logfmtBuffer() logfmt.Buffer {
	var buf logfmt.Buffer

	buf.WriteTimestamp(m.Timestamp)
	buf.WriteKey(m.Level.String())
	buf.WriteProperty("msg", m.Text)
	if m.Err != nil {
		buf.WriteProperty("error", m.Err.Error())
	}

	for _, p := range m.Properties {
		buf.WriteProperty(p.Key, p.Value)
	}

	for _, p := range m.Context {
		buf.WriteProperty(p.Key, p.Value)
	}

	if m.code != "" {
		buf.WriteProperty("code", m.code)
	}

	if m.status != 0 {
		buf.WriteProperty("status", m.status)
	}

	return buf
}
