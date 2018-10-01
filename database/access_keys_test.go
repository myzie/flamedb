package database

import (
	"testing"

	"github.com/myzie/flamedb/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessKey(t *testing.T) {

	key, text, err := NewAccessKey("my-key", database.ReadOnly)
	require.Nil(t, err)
	require.NotNil(t, key)

	assert.True(t, len(text) > 16)
	assert.NotNil(t, key.Compare("wrong"))
	assert.Nil(t, key.Compare(text))
}
