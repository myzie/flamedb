package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecord(t *testing.T) {
	r := &Record{}
	r.MustSetProperties(map[string]interface{}{"foo": "bar"})

	assert.Equal(t, map[string]interface{}{"foo": "bar"}, r.MustGetProperties())
}
