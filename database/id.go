package database

import (
	"io"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

var entropy io.Reader

func init() {
	t := time.Now()
	entropy = ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
}

// NewID returns a new unique identifier
func NewID() string {
	return ulid.MustNew(ulid.Now(), entropy).String()
}
