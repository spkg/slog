package slog

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	t1 := time.Now()
	m := newMessage(ctx, LevelDebug, "message text")
	t2 := time.Now()
	assert.False(m.Timestamp.Before(t1))
	assert.False(t2.Before(m.Timestamp))

	m.applyOpt(WithCode("xxx"))
	assert.Equal("xxx", m.code)

	errSample := errors.New("error here")
	m.applyOpts([]Option{WithStatus(1), WithValue("a", "b"), WithError(errSample)})
	assert.Equal(1, m.status)
	assert.Exactly(errSample, m.Err)
}

func TestMessageError(t *testing.T) {
	assert := assert.New(t)

	// set output to be discarded, and reset at end of test
	SetOutput(ioutil.Discard)
	defer func() { Default = New() }()

	ctx := context.Background()
	ctx = NewContext(ctx, Property{"a", "b"})
	m := Error(ctx, "This is the error",
		WithValue("c", "d"))
	assert.Equal("This is the error", m.Error())
}

func TestMessageCode(t *testing.T) {
	assert := assert.New(t)
	m := Message{}
	assert.Empty(m.Code())
	m.SetCode("xxx")
	assert.Equal("xxx", m.Code())
}

func TestMessageStatus(t *testing.T) {
	assert := assert.New(t)
	m := Message{}
	assert.Empty(m.Status())
	m.SetStatus(400)
	assert.Equal(400, m.Status())
}

func TestMessageLogfmt(t *testing.T) {
	assert := assert.New(t)
	m := Message{
		Timestamp: time.Unix(1234567890, 987654321).UTC(),
		Level:     LevelError,
		Text:      "This is the message",
		Err:       errors.New("Error message"),
		Properties: []Property{
			{"a", "b"},
			{"c", "d"},
		},
		Context: []Property{
			{"e", "f"},
			{"g", "h"},
		},
		code:   "CODE",
		status: 400,
	}

	expected := `2009-02-13T23:31:30.987654+0000 error msg="This is the message"` +
		` error="Error message" a=b c=d e=f g=h code=CODE status=400`
	assert.Equal(expected, m.Logfmt())
}
