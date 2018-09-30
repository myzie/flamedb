package database

import (
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AccessKey is a database entry allowing access to FlameDB
type AccessKey struct {
	ID        string    `json:"id" gorm:"size:64;primary_key;unique_index"`
	Name      string    `json:"name" gorm:"size:64"`
	CreatedAt time.Time `json:"created_at"`
	Value     string    `json:"value" gorm:"size:100"`
}

// Compare the given plaintext against this key to see if it's a match.
// Returns nil if it matches.
func (key *AccessKey) Compare(plaintext string) error {
	return bcryptCompare(key.Value, plaintext)
}

// NewAccessKey initializes a new AccessKey with a random secret value.
// The AccessKey model is returned along with the plaintext value which
// should be stored externally.
func NewAccessKey(name string) (*AccessKey, string, error) {

	value := NewID()

	cipher, err := bcryptHash(value)
	if err != nil {
		return nil, "", err
	}

	key := AccessKey{
		ID:        NewID(),
		Name:      name,
		CreatedAt: time.Now(),
		Value:     cipher,
	}
	return &key, value, nil
}

func bcryptHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}

func bcryptCompare(hashedPassword, password string) error {
	hash, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
