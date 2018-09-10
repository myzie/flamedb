package flamedb

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

var entropy *rand.Rand

func init() {
	entropy = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// NewID returns a new unique identifier
func NewID() string {
	return ulid.MustNew(ulid.Now(), entropy).String()
}
