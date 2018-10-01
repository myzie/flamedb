// Code generated by go-swagger; DO NOT EDIT.

package records

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/myzie/flamedb/models"
)

// NewCreateRecordParams creates a new CreateRecordParams object
// no default values defined in spec.
func NewCreateRecordParams() CreateRecordParams {

	return CreateRecordParams{}
}

// CreateRecordParams contains all the bound params for the create record operation
// typically these are obtained from a http.Request
//
// swagger:parameters createRecord
type CreateRecordParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Override user ID
	  In: header
	*/
	XUserID *string
	/*Record to be created
	  Required: true
	  In: body
	*/
	Body *models.RecordInput
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCreateRecordParams() beforehand.
func (o *CreateRecordParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := o.bindXUserID(r.Header[http.CanonicalHeaderKey("X-User-ID")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.RecordInput
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("body", "body"))
			} else {
				res = append(res, errors.NewParseError("body", "body", "", err))
			}
		} else {

			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Body = &body
			}
		}
	} else {
		res = append(res, errors.Required("body", "body"))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CreateRecordParams) bindXUserID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.XUserID = &raw

	return nil
}
