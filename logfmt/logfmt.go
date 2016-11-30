// Package logfmt provides some helper functions to write logfmt
// messages. See https://brandur.org/logfmt for a description of
// the logfmt format.
package logfmt

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"runtime"
	"time"
)

var (
	// TimeFormat is the time format used for timestamps.
	TimeFormat = "2006-01-02T15:04:05.000000-0700"
)

var (
	// end of line bytes
	eol []byte
)

func init() {
	if runtime.GOOS == "windows" {
		eol = []byte{0xd, 0x0a}
	} else {
		eol = []byte{0x0a}
	}
}

// Buffer is used for building up logfmt messages. Once the
// message is built, the text can be obtained by the String method,
// or the message can be written to an io.Writer using the WriteTo method.
// The calling program should call the Reset method after it is finished with
// the buffer because then internal buffers can be re-used to take pressure off
// the garbage collector.
type Buffer struct {
	buf *bytes.Buffer
}

// Reset releases internal memory used by the buffer and makes that memory
// available for reuse by other Buffers. Calling Reset in programs that log many
// messages will reduce the pressure on the garbage collector.
func (b *Buffer) Reset() {
	if b.buf != nil {
		releaseBuffer(b.buf)
		b.buf = nil
	}
}

// String returns a string representation of the message in the buffer.
// It implements the fmt.Stringer interface.
func (b *Buffer) String() string {
	b.allocate()
	return string(b.buf.Bytes())
}

// Len returns the length of the message stored in the buffer.
func (b *Buffer) Len() int {
	b.allocate()
	return b.buf.Len()
}

// WriteTo implements the io.WriterTo interface. It writes the formatted
// message to the writer w.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	b.allocate()
	n, err := w.Write(b.buf.Bytes())
	return int64(n), err
}

// WriteTimestamp writes a timestamp to the buffer. The format of the timestamp
// is determined by the TimestampFormat variable.
func (b *Buffer) WriteTimestamp(t time.Time) error {
	b.allocate()
	if err := b.spacer(); err != nil {
		return err
	}
	_, err := b.buf.WriteString(t.Format(TimeFormat))
	return err
}

// WriteKey writes a single key without a value to the buffer. The key should not
// contain any special characters.
func (b *Buffer) WriteKey(key string) error {
	b.allocate()
	if err := b.spacer(); err != nil {
		return err
	}
	_, err := b.buf.WriteString(key)
	return err
}

// WriteProperty writes a key value pair to the buffer. If the value contains any special
// characters it will be quoted.
func (b *Buffer) WriteProperty(key string, value interface{}) error {
	b.allocate()
	if err := b.spacer(); err != nil {
		return err
	}
	return writeProperty(b.buf, key, value)
}

// allocate ensures that a buffer is allocated.
func (b *Buffer) allocate() {
	if b.buf == nil {
		b.buf = getBuffer()
	}
}

// spacer adds a space to the buffer if it is not empty.
// Called by functions that are going to add more to the buffer.
func (b *Buffer) spacer() error {
	if b.buf.Len() > 0 {
		_, err := b.buf.WriteRune(' ')
		return err
	}
	return nil
}

// WriteNewLine writes a single new line byte to the buffer.
func (b *Buffer) WriteNewLine() error {
	b.allocate()
	_, err := b.buf.Write(eol)
	return err
}

// WriteEOL writes the OS-specific new line bytes to the buffer.
// On Windows the new line is 0xd, 0xa. For all other operating systems
// the new line is 0x0a.
func (b *Buffer) WriteEOL() error {
	b.allocate()
	_, err := b.buf.Write(eol)
	return err
}

// writeValueString writes a string value to the buffer buf.
// If the contains any characters that need quoting, then the
// value is written in quotes. Otherwise the value is written
// without quotes.
func writeValueString(buf *bytes.Buffer, value string) error {
	var err error
	needsQuotes := false
	needsEscape := false
	hasEscape := false
	for _, c := range value {
		if c == '"' || c == '=' || c == '\r' || c == '\n' || c == '\t' {
			needsQuotes = true
			needsEscape = true
		} else if c <= ' ' {
			needsQuotes = true
		} else if c == '\\' {
			hasEscape = true
		}
	}

	if needsQuotes && hasEscape {
		needsEscape = true
	}

	if needsQuotes {
		_, err = buf.WriteRune('"')
		if err != nil {
			return err
		}
		if needsEscape {
			// need to escape one or more chars in the value
			for _, c := range value {
				switch c {
				case '\r':
					// ignore CR in message
				case '\n':
					_, err = buf.WriteString("\\n")
					if err != nil {
						return err
					}
				case '\t':
					_, err = buf.WriteString("\\t")
					if err != nil {
						return err
					}
				default:
					if c == '\\' || c == '"' {
						_, err = buf.WriteRune('\\')
						if err != nil {
							return err
						}
					}
					_, err = buf.WriteRune(c)
					if err != nil {
						return err
					}
				}
			}
		} else {
			// no escapable chars in the value, so can just write the contents
			_, err := buf.WriteString(value)
			if err != nil {
				return err
			}
		}
		_, err := buf.WriteRune('"')
		if err != nil {
			return err
		}
	} else {
		// no quotes required, just write the string
		_, err := buf.WriteString(value)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeProperty writes a key value pair to buf in a format compatible with logfmt.
func writeProperty(buf *bytes.Buffer, key string, value interface{}) error {
	var err error
	_, err = buf.WriteString(key)
	if err != nil {
		return err
	}
	_, err = buf.WriteRune('=')
	if err != nil {
		return err
	}
	switch v := value.(type) {
	case bool:
		_, err = fmt.Fprint(buf, v)
		return err
	case byte:
		_, err = fmt.Fprint(buf, v)
		return err
	case complex64:
		return writeValueString(buf, fmt.Sprint(v))
	case complex128:
		return writeValueString(buf, fmt.Sprint(v))
	case error:
		return writeValueString(buf, v.Error())
	case float32:
		_, err = fmt.Fprint(buf, v)
		return err
	case float64:
		_, err = fmt.Fprint(buf, v)
		return err
	case int:
		_, err = fmt.Fprint(buf, v)
		return err
	case int16:
		_, err = fmt.Fprint(buf, v)
		return err
	case int32:
		_, err = fmt.Fprint(buf, v)
		return err
	case int64:
		_, err = fmt.Fprint(buf, v)
		return err
	case int8:
		_, err = fmt.Fprint(buf, v)
		return err
	case string:
		return writeValueString(buf, v)
	case time.Time:
		return writeValueString(buf, v.Format(TimeFormat))
	case uint:
		_, err = fmt.Fprint(buf, v)
		return err
	case uint16:
		_, err = fmt.Fprint(buf, v)
		return err
	case uint32:
		_, err = fmt.Fprint(buf, v)
		return err
	case uint64:
		_, err = fmt.Fprint(buf, v)
		return err
	case uintptr:
		_, err = fmt.Fprint(buf, v)
		return err
	}

	if v, ok := value.(fmt.Stringer); ok {
		return writeValueString(buf, v.String())
	}

	if v, ok := value.(encoding.TextMarshaler); ok {
		b, err := v.MarshalText()
		if err != nil {
			return err
		}
		return writeValueString(buf, string(b))
	}

	return writeValueString(buf, fmt.Sprint(value))
}
