// Code generated by go-swagger; DO NOT EDIT.

package records

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/myzie/flamedb/models"
)

// UpdateRecordHandlerFunc turns a function with the right signature into a update record handler
type UpdateRecordHandlerFunc func(UpdateRecordParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateRecordHandlerFunc) Handle(params UpdateRecordParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// UpdateRecordHandler interface for that can handle valid update record params
type UpdateRecordHandler interface {
	Handle(UpdateRecordParams, *models.Principal) middleware.Responder
}

// NewUpdateRecord creates a new http.Handler for the update record operation
func NewUpdateRecord(ctx *middleware.Context, handler UpdateRecordHandler) *UpdateRecord {
	return &UpdateRecord{Context: ctx, Handler: handler}
}

/*UpdateRecord swagger:route PUT /records/{recordId} records updateRecord

Update an existing record

*/
type UpdateRecord struct {
	Context *middleware.Context
	Handler UpdateRecordHandler
}

func (o *UpdateRecord) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateRecordParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
