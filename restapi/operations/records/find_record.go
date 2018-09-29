// Code generated by go-swagger; DO NOT EDIT.

package records

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/myzie/flamedb/models"
)

// FindRecordHandlerFunc turns a function with the right signature into a find record handler
type FindRecordHandlerFunc func(FindRecordParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn FindRecordHandlerFunc) Handle(params FindRecordParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// FindRecordHandler interface for that can handle valid find record params
type FindRecordHandler interface {
	Handle(FindRecordParams, *models.Principal) middleware.Responder
}

// NewFindRecord creates a new http.Handler for the find record operation
func NewFindRecord(ctx *middleware.Context, handler FindRecordHandler) *FindRecord {
	return &FindRecord{Context: ctx, Handler: handler}
}

/*FindRecord swagger:route GET /find records findRecord

Find a record by path

*/
type FindRecord struct {
	Context *middleware.Context
	Handler FindRecordHandler
}

func (o *FindRecord) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewFindRecordParams()

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
