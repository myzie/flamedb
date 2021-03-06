// Code generated by go-swagger; DO NOT EDIT.

package records

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/myzie/flamedb/models"
)

// CreateRecordOKCode is the HTTP code returned for type CreateRecordOK
const CreateRecordOKCode int = 200

/*CreateRecordOK Successfully created

swagger:response createRecordOK
*/
type CreateRecordOK struct {

	/*
	  In: Body
	*/
	Payload *models.RecordOutput `json:"body,omitempty"`
}

// NewCreateRecordOK creates CreateRecordOK with default headers values
func NewCreateRecordOK() *CreateRecordOK {

	return &CreateRecordOK{}
}

// WithPayload adds the payload to the create record o k response
func (o *CreateRecordOK) WithPayload(payload *models.RecordOutput) *CreateRecordOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create record o k response
func (o *CreateRecordOK) SetPayload(payload *models.RecordOutput) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateRecordOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateRecordBadRequestCode is the HTTP code returned for type CreateRecordBadRequest
const CreateRecordBadRequestCode int = 400

/*CreateRecordBadRequest Bad request

swagger:response createRecordBadRequest
*/
type CreateRecordBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.BadRequest `json:"body,omitempty"`
}

// NewCreateRecordBadRequest creates CreateRecordBadRequest with default headers values
func NewCreateRecordBadRequest() *CreateRecordBadRequest {

	return &CreateRecordBadRequest{}
}

// WithPayload adds the payload to the create record bad request response
func (o *CreateRecordBadRequest) WithPayload(payload *models.BadRequest) *CreateRecordBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create record bad request response
func (o *CreateRecordBadRequest) SetPayload(payload *models.BadRequest) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateRecordBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateRecordInternalServerErrorCode is the HTTP code returned for type CreateRecordInternalServerError
const CreateRecordInternalServerErrorCode int = 500

/*CreateRecordInternalServerError Internal Server Error

swagger:response createRecordInternalServerError
*/
type CreateRecordInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.InternalServerError `json:"body,omitempty"`
}

// NewCreateRecordInternalServerError creates CreateRecordInternalServerError with default headers values
func NewCreateRecordInternalServerError() *CreateRecordInternalServerError {

	return &CreateRecordInternalServerError{}
}

// WithPayload adds the payload to the create record internal server error response
func (o *CreateRecordInternalServerError) WithPayload(payload *models.InternalServerError) *CreateRecordInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create record internal server error response
func (o *CreateRecordInternalServerError) SetPayload(payload *models.InternalServerError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateRecordInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
