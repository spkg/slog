package logfmt

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type stringer int

func (s stringer) String() string {
	return fmt.Sprintf("stringer: %d", s)
}

type textMarshaler int

func (tm textMarshaler) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("textMarshaler: %d", tm)), nil
}

type needsSprintf struct {
	A int
	B string
}

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
		{Value: byte(0x10), Expected: "key=16"},
		{Value: complex(float32(10), float32(11)), Expected: "key=(10+11i)"},
		{Value: complex(float64(10.4), float64(11.5)), Expected: "key=(10.4+11.5i)"},
		{Value: errors.New("This is an error"), Expected: `key="This is an error"`},
		{Value: float32(3.14159), Expected: "key=3.14159"},
		{Value: float64(31.4159), Expected: "key=31.4159"},
		{Value: int(1), Expected: "key=1"},
		{Value: int16(2), Expected: "key=2"},
		{Value: int32(3), Expected: "key=3"},
		{Value: int64(4), Expected: "key=4"},
		{Value: int8(5), Expected: "key=5"},
		{Value: "string", Expected: `key=string`},
		{Value: uint(1), Expected: "key=1"},
		{Value: uint16(2), Expected: "key=2"},
		{Value: uint32(3), Expected: "key=3"},
		{Value: uint64(4), Expected: "key=4"},
		{Value: uintptr(3041255), Expected: "key=3041255"},
		{Value: stringer(44), Expected: `key="stringer: 44"`},
		{Value: textMarshaler(45), Expected: `key="textMarshaler: 45"`},
		{Value: needsSprintf{46, "text value"}, Expected: `key="{46 text value}"`},
	}

	for _, tc := range testCases {
		buf := Buffer{}
		buf.WriteProperty("key", tc.Value)
		assert.Equal(tc.Expected, buf.String())
	}
}
