// Code generated by go-swagger; DO NOT EDIT.

package records

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/myzie/flamedb/models"
)

// FindRecordOKCode is the HTTP code returned for type FindRecordOK
const FindRecordOKCode int = 200

/*FindRecordOK The retrieved record

swagger:response findRecordOK
*/
type FindRecordOK struct {

	/*
	  In: Body
	*/
	Payload *models.RecordOutput `json:"body,omitempty"`
}

// NewFindRecordOK creates FindRecordOK with default headers values
func NewFindRecordOK() *FindRecordOK {

	return &FindRecordOK{}
}

// WithPayload adds the payload to the find record o k response
func (o *FindRecordOK) WithPayload(payload *models.RecordOutput) *FindRecordOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find record o k response
func (o *FindRecordOK) SetPayload(payload *models.RecordOutput) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindRecordOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// FindRecordNotFoundCode is the HTTP code returned for type FindRecordNotFound
const FindRecordNotFoundCode int = 404

/*FindRecordNotFound Record not found

swagger:response findRecordNotFound
*/
type FindRecordNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.NotFoundError `json:"body,omitempty"`
}

// NewFindRecordNotFound creates FindRecordNotFound with default headers values
func NewFindRecordNotFound() *FindRecordNotFound {

	return &FindRecordNotFound{}
}

// WithPayload adds the payload to the find record not found response
func (o *FindRecordNotFound) WithPayload(payload *models.NotFoundError) *FindRecordNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find record not found response
func (o *FindRecordNotFound) SetPayload(payload *models.NotFoundError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindRecordNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// FindRecordInternalServerErrorCode is the HTTP code returned for type FindRecordInternalServerError
const FindRecordInternalServerErrorCode int = 500

/*FindRecordInternalServerError Internal Server Error

swagger:response findRecordInternalServerError
*/
type FindRecordInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.InternalServerError `json:"body,omitempty"`
}

// NewFindRecordInternalServerError creates FindRecordInternalServerError with default headers values
func NewFindRecordInternalServerError() *FindRecordInternalServerError {

	return &FindRecordInternalServerError{}
}

// WithPayload adds the payload to the find record internal server error response
func (o *FindRecordInternalServerError) WithPayload(payload *models.InternalServerError) *FindRecordInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find record internal server error response
func (o *FindRecordInternalServerError) SetPayload(payload *models.InternalServerError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindRecordInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
