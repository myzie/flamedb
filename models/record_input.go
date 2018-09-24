// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// RecordInput record input
// swagger:model RecordInput
type RecordInput struct {

	// id
	ID string `json:"id,omitempty"`

	// path
	// Required: true
	Path *string `json:"path"`

	// properties
	// Required: true
	Properties map[string]string `json:"properties"`
}

// Validate validates this record input
func (m *RecordInput) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePath(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateProperties(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RecordInput) validatePath(formats strfmt.Registry) error {

	if err := validate.Required("path", "body", m.Path); err != nil {
		return err
	}

	return nil
}

func (m *RecordInput) validateProperties(formats strfmt.Registry) error {

	return nil
}

// MarshalBinary interface implementation
func (m *RecordInput) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RecordInput) UnmarshalBinary(b []byte) error {
	var res RecordInput
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
