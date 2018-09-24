// Code generated by go-swagger; DO NOT EDIT.

package records

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/myzie/flamedb/models"
)

// GetRecordOKCode is the HTTP code returned for type GetRecordOK
const GetRecordOKCode int = 200

/*GetRecordOK The retrieved record

swagger:response getRecordOK
*/
type GetRecordOK struct {

	/*
	  In: Body
	*/
	Payload *models.RecordOutput `json:"body,omitempty"`
}

// NewGetRecordOK creates GetRecordOK with default headers values
func NewGetRecordOK() *GetRecordOK {

	return &GetRecordOK{}
}

// WithPayload adds the payload to the get record o k response
func (o *GetRecordOK) WithPayload(payload *models.RecordOutput) *GetRecordOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get record o k response
func (o *GetRecordOK) SetPayload(payload *models.RecordOutput) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRecordOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetRecordNotFoundCode is the HTTP code returned for type GetRecordNotFound
const GetRecordNotFoundCode int = 404

/*GetRecordNotFound Record not found

swagger:response getRecordNotFound
*/
type GetRecordNotFound struct {
}

// NewGetRecordNotFound creates GetRecordNotFound with default headers values
func NewGetRecordNotFound() *GetRecordNotFound {

	return &GetRecordNotFound{}
}

// WriteResponse to the client
func (o *GetRecordNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}