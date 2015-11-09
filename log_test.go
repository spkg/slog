package slog

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func TestLog(t *testing.T) {
	assert := assert.New(t)
	Default = New()
	defer func() { Default = New() }()
	SetOutput(ioutil.Discard)
	th := &testHandler{}
	AddHandler(th)
	ctx := context.Background()

	// should only log 3 messages: debug not enabled by default
	Debug(ctx, "debug message")
	Info(ctx, "info message")
	Warn(ctx, "warn message")
	Error(ctx, "error message")
	assert.Equal(3, len(th.Messages))

	SetMinLevel(LevelDebug)
	th.Messages = nil

	// should log 4 messages: debug now enabled
	Debug(ctx, "debug message")
	Info(ctx, "info message")
	Warn(ctx, "warn message")
	Error(ctx, "error message")
	assert.Equal(4, len(th.Messages))

}

type testHandler struct {
	Messages []*Message
}

func (th *testHandler) Handle(msgs []*Message) {
	th.Messages = append(th.Messages, msgs...)
}
