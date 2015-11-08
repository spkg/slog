package logfmt

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriteTo(t *testing.T) {
	assert := assert.New(t)

	buf := Buffer{}
	tm := time.Unix(1234567890, 987654321).UTC()
	buf.WriteTimestamp(tm)
	buf.WriteKey("info")
	buf.WriteProperty("key", "value")

	bytes := bytes.Buffer{}
	buf.WriteTo(&bytes)
	text := string(bytes.Bytes())
	assert.Equal("2009-02-13T23:31:30.987654+0000 info key=value", text)
}

func TestFormatting(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Value    interface{}
		Expected string
	}{
		{Value: true, Expected: "key=true"},
		{Value: false, Expected: "key=false"},
		{Value: complex(float32(10), float32(11)), Expected: "key=(10+11i)"},
		{Value: complex(float64(10.4), float64(11.5)), Expected: "key=(10.4+11.5i)"},
		{Value: byte(0x10), Expected: "key=16"},
		{Value: int(1), Expected: "key=1"},
	}

	for _, tc := range testCases {
		buf := Buffer{}
		buf.WriteProperty("key", tc.Value)
		assert.Equal(tc.Expected, buf.String())
	}
}
