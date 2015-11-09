package slog

import (
	"errors"
	"fmt"
	"strings"
)

// Level indicates the level of a log message.
type Level int

const (
	LevelDebug   Level = iota // Debugging only
	LevelInfo                 // Informational
	LevelWarning              // Warning
	LevelError                // Error condition
)

// MinLevel is the minimum level that will be logged.
// The calling program can change this value at any time.
var (
	MinLevel = LevelInfo
)

var (
	errInvalidLevel = errors.New("invalid level")
)

// String implements the String interface.
func (s Level) String() string {
	switch s {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warn"
	case LevelError:
		return "error"
	}
	return fmt.Sprintf("unknown %d", s)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (lvl Level) MarshalText() ([]byte, error) {
	return []byte(lvl.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (lvl *Level) UnmarshalText(text []byte) error {
	str := strings.ToLower(string(text))
	switch str {
	case "debug":
		*lvl = LevelDebug
	case "info", "information":
		*lvl = LevelInfo
	case "warn", "warning":
		*lvl = LevelWarning
	case "error":
		*lvl = LevelError
	default:
		return errInvalidLevel
	}
	return nil
}
