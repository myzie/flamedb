// Code generated by go-swagger; DO NOT EDIT.

package records

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewListRecordsParams creates a new ListRecordsParams object
// no default values defined in spec.
func NewListRecordsParams() ListRecordsParams {

	return ListRecordsParams{}
}

// ListRecordsParams contains all the bound params for the list records operation
// typically these are obtained from a http.Request
//
// swagger:parameters listRecords
type ListRecordsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Override user ID
	  In: header
	*/
	XUserID *string
	/*
	  Maximum: 1000
	  Minimum: 0
	  In: query
	*/
	Limit *int64
	/*
	  Minimum: 0
	  In: query
	*/
	Offset *int64
	/*
	  In: query
	*/
	OrderBy *string
	/*
	  In: query
	*/
	OrderByDesc *bool
	/*
	  In: query
	*/
	OrderByProperty *string
	/*
	  In: query
	*/
	OrderByPropertyDesc *bool
	/*
	  In: query
	*/
	Parent *string
	/*
	  In: query
	*/
	Prefix *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewListRecordsParams() beforehand.
func (o *ListRecordsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	if err := o.bindXUserID(r.Header[http.CanonicalHeaderKey("X-User-ID")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	qLimit, qhkLimit, _ := qs.GetOK("limit")
	if err := o.bindLimit(qLimit, qhkLimit, route.Formats); err != nil {
		res = append(res, err)
	}

	qOffset, qhkOffset, _ := qs.GetOK("offset")
	if err := o.bindOffset(qOffset, qhkOffset, route.Formats); err != nil {
		res = append(res, err)
	}

	qOrderBy, qhkOrderBy, _ := qs.GetOK("orderBy")
	if err := o.bindOrderBy(qOrderBy, qhkOrderBy, route.Formats); err != nil {
		res = append(res, err)
	}

	qOrderByDesc, qhkOrderByDesc, _ := qs.GetOK("orderByDesc")
	if err := o.bindOrderByDesc(qOrderByDesc, qhkOrderByDesc, route.Formats); err != nil {
		res = append(res, err)
	}

	qOrderByProperty, qhkOrderByProperty, _ := qs.GetOK("orderByProperty")
	if err := o.bindOrderByProperty(qOrderByProperty, qhkOrderByProperty, route.Formats); err != nil {
		res = append(res, err)
	}

	qOrderByPropertyDesc, qhkOrderByPropertyDesc, _ := qs.GetOK("orderByPropertyDesc")
	if err := o.bindOrderByPropertyDesc(qOrderByPropertyDesc, qhkOrderByPropertyDesc, route.Formats); err != nil {
		res = append(res, err)
	}

	qParent, qhkParent, _ := qs.GetOK("parent")
	if err := o.bindParent(qParent, qhkParent, route.Formats); err != nil {
		res = append(res, err)
	}

	qPrefix, qhkPrefix, _ := qs.GetOK("prefix")
	if err := o.bindPrefix(qPrefix, qhkPrefix, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ListRecordsParams) bindXUserID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

func (o *ListRecordsParams) bindLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("limit", "query", "int64", raw)
	}
	o.Limit = &value

	if err := o.validateLimit(formats); err != nil {
		return err
	}

	return nil
}

func (o *ListRecordsParams) validateLimit(formats strfmt.Registry) error {

	if err := validate.MinimumInt("limit", "query", int64(*o.Limit), 0, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("limit", "query", int64(*o.Limit), 1000, false); err != nil {
		return err
	}

	return nil
}

func (o *ListRecordsParams) bindOffset(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("offset", "query", "int64", raw)
	}
	o.Offset = &value

	if err := o.validateOffset(formats); err != nil {
		return err
	}

	return nil
}

func (o *ListRecordsParams) validateOffset(formats strfmt.Registry) error {

	if err := validate.MinimumInt("offset", "query", int64(*o.Offset), 0, false); err != nil {
		return err
	}

	return nil
}

func (o *ListRecordsParams) bindOrderBy(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.OrderBy = &raw

	return nil
}

func (o *ListRecordsParams) bindOrderByDesc(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("orderByDesc", "query", "bool", raw)
	}
	o.OrderByDesc = &value

	return nil
}

func (o *ListRecordsParams) bindOrderByProperty(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.OrderByProperty = &raw

	return nil
}

func (o *ListRecordsParams) bindOrderByPropertyDesc(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("orderByPropertyDesc", "query", "bool", raw)
	}
	o.OrderByPropertyDesc = &value

	return nil
}

func (o *ListRecordsParams) bindParent(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Parent = &raw

	return nil
}

func (o *ListRecordsParams) bindPrefix(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Prefix = &raw

	return nil
}
