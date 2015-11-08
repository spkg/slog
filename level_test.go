package slog_test

import (
	"encoding/json"
	"testing"

	"sp.com.au/exp/log"

	"github.com/stretchr/testify/assert"
)

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
