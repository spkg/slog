package slog_test

import (
	"encoding/json"
	"testing"

	"sp.com.au/exp/log"

	"github.com/spkg/slog"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Level    slog.Level
		Expected string
	}{
		{slog.LevelDebug, "debug"},
		{slog.LevelInfo, "info"},
		{slog.LevelWarning, "warn"},
		{slog.LevelError, "error"},
		{slog.Level(63), "unknown 63"},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, tc.Level.String())
	}
}

func TestMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	s := struct{ Level log.Level }{log.LevelError}
	b, err := json.Marshal(&s)
	assert.NoError(err)
	assert.Equal(`{"Level":"error"}`, string(b))

	var s2 struct{ Level log.Level }
	err = json.Unmarshal(b, &s2)
	assert.NoError(err)
	assert.Equal(s.Level, s2.Level)
}

func TestMarshalText(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Level    slog.Level
		Expected string
		Valid    []string
		Invalid  []string
	}{
		{Level: slog.LevelDebug, Expected: "debug", Valid: []string{"Debug", "DEBUG"}},
		{Level: slog.LevelInfo, Expected: "info", Valid: []string{"INFO", "information"}, Invalid: []string{"xxxx"}},
		{Level: slog.LevelWarning, Expected: "warn", Valid: []string{"Warning", "WARN"}},
		{Level: slog.LevelError, Expected: "error"},
	}

	for _, tc := range testCases {
		b, err := tc.Level.MarshalText()
		assert.NoError(err)
		assert.Equal(tc.Expected, string(b))

		var level slog.Level
		err = level.UnmarshalText(b)
		assert.NoError(err)
		assert.Equal(tc.Level, level)

		for _, text := range tc.Valid {
			err = level.UnmarshalText([]byte(text))
			assert.NoError(err)
			assert.Equal(tc.Level, level)
		}
		for _, text := range tc.Invalid {
			err = level.UnmarshalText([]byte(text))
			assert.Error(err)
		}
	}
}
