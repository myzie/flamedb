// Code generated by go-swagger; DO NOT EDIT.

package records

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/myzie/flamedb/models"
)

// ListRecordsOKCode is the HTTP code returned for type ListRecordsOK
const ListRecordsOKCode int = 200

/*ListRecordsOK Successful query

swagger:response listRecordsOK
*/
type ListRecordsOK struct {

	/*
	  In: Body
	*/
	Payload *models.QueryResult `json:"body,omitempty"`
}

// NewListRecordsOK creates ListRecordsOK with default headers values
func NewListRecordsOK() *ListRecordsOK {

	return &ListRecordsOK{}
}

// WithPayload adds the payload to the list records o k response
func (o *ListRecordsOK) WithPayload(payload *models.QueryResult) *ListRecordsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list records o k response
func (o *ListRecordsOK) SetPayload(payload *models.QueryResult) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListRecordsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
