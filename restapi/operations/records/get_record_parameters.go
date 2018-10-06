// Code generated by go-swagger; DO NOT EDIT.

package records

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetRecordParams creates a new GetRecordParams object
// no default values defined in spec.
func NewGetRecordParams() GetRecordParams {

	return GetRecordParams{}
}

// GetRecordParams contains all the bound params for the get record operation
// typically these are obtained from a http.Request
//
// swagger:parameters getRecord
type GetRecordParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Override user ID
	  In: header
	*/
	XUserID *string
	/*ID of record to return
	  Required: true
	  In: path
	*/
	RecordID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetRecordParams() beforehand.
func (o *GetRecordParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := o.bindXUserID(r.Header[http.CanonicalHeaderKey("X-User-ID")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	rRecordID, rhkRecordID, _ := route.Params.GetOK("recordId")
	if err := o.bindRecordID(rRecordID, rhkRecordID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetRecordParams) bindXUserID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

func (o *GetRecordParams) bindRecordID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.RecordID = raw

	return nil
}
