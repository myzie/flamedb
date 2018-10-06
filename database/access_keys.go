package database

//go:generate mockgen -source access_keys.go -destination mock_database/mock_access_key_store.go

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	gorm "github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// AccessKeyPermission defines the access level associated with an access key
type AccessKeyPermission string

const (
	// Denied access
	Denied AccessKeyPermission = ""

	// ReadOnly grant
	ReadOnly AccessKeyPermission = "r"

	// ReadWrite grant
	ReadWrite AccessKeyPermission = "rw"

	// ServiceRead provides read access from an external service
	ServiceRead AccessKeyPermission = "sr"

	// ServiceReadWrite provides rw access from an external service
	ServiceReadWrite AccessKeyPermission = "srw"
)

// AccessKey is a database entry allowing access to FlameDB
type AccessKey struct {
	ID         string    `json:"id" gorm:"size:64;primary_key;unique_index"`
	RefID      string    `json:"ref_id" gorm:"size:64;index"`
	Name       string    `json:"name" gorm:"size:64"`
	Permission string    `json:"permission" gorm:"size:64"`
	CreatedAt  time.Time `json:"created_at"`
	Secret     string    `json:"value" gorm:"size:100"`
}

// Compare the given plaintext against this key to see if it's a match.
// Returns nil if it matches.
func (key *AccessKey) Compare(plaintext string) error {
	return bcryptCompare(key.Secret, plaintext)
}

// NewAccessKey initializes a new AccessKey with a random secret value.
// The AccessKey model is returned along with the plaintext value which
// should be stored externally.
func NewAccessKey(name, refID string, perm AccessKeyPermission) (*AccessKey, string, error) {

	switch perm {
	case ReadOnly:
	case ReadWrite:
	case ServiceRead:
	case ServiceReadWrite:
	default:
		return nil, "", errors.New("Invalid permission value")
	}

	keyID, err := uuid.NewV4()
	if err != nil {
		return nil, "", err
	}

	keySecret, err := uuid.NewV4()
	if err != nil {
		return nil, "", err
	}

	cipherText, err := bcryptHash(keySecret.String())
	if err != nil {
		return nil, "", err
	}

	key := AccessKey{
		ID:         keyID.String(),
		RefID:      refID,
		Name:       name,
		Permission: string(perm),
		CreatedAt:  time.Now(),
		Secret:     cipherText,
	}
	return &key, keySecret.String(), nil
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

// AccessKeyStore is an interface to retrieve access permissions
type AccessKeyStore interface {

	// Get returns the access key with the given key ID
	Get(id string) (*AccessKey, error)
}

// NewAccessKeyStore returns an initialized access key store
func NewAccessKeyStore(gormDB *gorm.DB) (AccessKeyStore, error) {

	if err := gormDB.AutoMigrate(&AccessKey{}).Error; err != nil {
		return nil, fmt.Errorf("Failed to migrate: %s", err.Error())
	}

	return &accessKeyStore{gormDB: gormDB}, nil
}

type accessKeyStore struct {
	gormDB *gorm.DB
}

func (ks *accessKeyStore) Get(id string) (*AccessKey, error) {
	// Retrieve the access key with that ID
	key := &AccessKey{}
	if err := ks.gormDB.Where(AccessKey{ID: id}).First(&key).Error; err != nil {
		return nil, err
	}
	return key, nil
}
