package database

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// Record is an item with a unique ID stored in FlameDB at a specified Path.
// The Record has a handful of top-level meta attributes and a Properties
// field which may contain arbitrary JSON.
type Record struct {
	ID         string         `json:"id" gorm:"size:64;primary_key;unique_index"`
	Path       string         `json:"path" gorm:"size:320;unique_index"`
	Parent     string         `json:"parent" gorm:"size:256;index"`
	CreatedBy  string         `json:"created_by" gorm:"size:64"`
	UpdatedBy  string         `json:"updated_by" gorm:"size:64"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	Properties postgres.Jsonb `json:"properties"`
}

// GetProperties returns the record properties as a map
func (r *Record) GetProperties() (map[string]interface{}, error) {
	var props map[string]interface{}
	if err := json.Unmarshal(r.Properties.RawMessage, &props); err != nil {
		return nil, err
	}
	return props, nil
}

// MustGetProperties returns the record properties as a map.
// If they fail to unmarshal this panics.
func (r *Record) MustGetProperties() map[string]interface{} {
	props, err := r.GetProperties()
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal properties: %s", err.Error()))
	}
	return props
}

// SetProperties sets the properites in JSONB format
func (r *Record) SetProperties(properties map[string]interface{}) error {
	propJSON, err := json.Marshal(properties)
	if err != nil {
		return err
	}
	r.Properties = postgres.Jsonb{RawMessage: json.RawMessage(propJSON)}
	return nil
}

// MustSetProperties sets properties and panics if marshaling fails
func (r *Record) MustSetProperties(properties map[string]interface{}) {
	if err := r.SetProperties(properties); err != nil {
		panic(err)
	}
}

// BeforeSave is called to validate the Record before saving it to the database
func (r *Record) BeforeSave() error {

	if r.ID == "" {
		return fmt.Errorf("Invalid record id: empty")
	}

	if len(r.Path) == 0 {
		return fmt.Errorf("Invalid record path: empty")
	} else if len(r.Path) > 320 {
		return fmt.Errorf("Invalid record path: too long")
	} else if r.Path[0] != '/' {
		return fmt.Errorf("Invalid record path: does not start with /")
	}

	parent, _ := path.Split(r.Path)
	r.Parent = parent

	if len(r.Properties.RawMessage) > 1024 {
		return fmt.Errorf("Invalid record properties: too large")
	}
	return nil
}
