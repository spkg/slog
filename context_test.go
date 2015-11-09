package slog

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	ctx = NewContext(ctx, Property{"property1", "a1"}, Property{"property2", 2})
	ctx = NewContext(ctx, Property{"property3", 12.34})

	logData := fromContext(ctx)
	assert.Equal("property3", logData.Key)
	assert.Equal(12.34, logData.Value)
	logData = logData.Prev
	assert.Equal("property1", logData.Key)
	assert.Equal("a1", logData.Value)
	logData = logData.Prev
	assert.Equal("property2", logData.Key)
	assert.Equal(2, logData.Value)
	assert.Nil(logData.Prev)
}

func TestNewContextNoProperties(t *testing.T) {
	assert := assert.New(t)
	ctx1 := context.Background()
	ctx2 := NewContext(ctx1)
	assert.Exactly(ctx1, ctx2)
	ctx2 = NewContext(ctx1, Property{"a", "b"})
	assert.NotEqual(ctx2, ctx1)
}

func TestNewContextNilContext(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(fromContext(nil))
}

func TestFromContext(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(fromContext(context.Background()))
}
