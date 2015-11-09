package slog

import (
	"bytes"
	"regexp"

	"golang.org/x/net/context"
)

var errorRegexp = regexp.MustCompile(`(?i)error|panic|fatal`)

// Writer implements the io.Writer interface and is intended to be supplied
// as the writer for a standard library logger.
type writer struct {
	ctx    context.Context
	logger Logger
}

func (w *writer) Write(b []byte) (int, error) {
	buf := bytes.NewBuffer(b)

	// Read the first line of text only. TODO: when there are multiple lines what to do with the rest.
	// Not checking for error, because the only error is that there is no terminating
	// line feed, which does not matter.
	text, _ := buf.ReadBytes(0x0a)

	if errorRegexp.Match(text) {
		w.logger.Error(w.ctx, string(text))
	} else {
		w.logger.Info(w.ctx, string(text))
	}

	return len(b), nil
}
