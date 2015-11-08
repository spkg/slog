package slog

import (
	"bytes"
	"strings"
)

// Writer implements the io.Writer interface and is intended to be supplied
// as the writer for a standard library logger.
type Writer struct{}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Write(b []byte) (int, error) {
	buf := bytes.NewBuffer(b)

	// Read the first line of text only. TODO: when there are multiple lines what to do with the rest.
	// Not checking for error, because the only error is that there is no terminating
	// line feed, which does not matter.
	text, _ := buf.ReadString(0x0a)
	lowerText := strings.ToLower(text)

	level := LevelWarning

	for _, m := range []string{"error", "panic", "fatal"} {
		if strings.Contains(lowerText, m) {
			level = LevelError
			break
		}
	}

	if level == LevelWarning {
		Warn(text)
	} else {
		Error(text)
	}

	return len(b), nil
}
