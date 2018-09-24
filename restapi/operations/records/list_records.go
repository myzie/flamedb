// Code generated by go-swagger; DO NOT EDIT.

package records

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// ListRecordsHandlerFunc turns a function with the right signature into a list records handler
type ListRecordsHandlerFunc func(ListRecordsParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn ListRecordsHandlerFunc) Handle(params ListRecordsParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// ListRecordsHandler interface for that can handle valid list records params
type ListRecordsHandler interface {
	Handle(ListRecordsParams, interface{}) middleware.Responder
}

// NewListRecords creates a new http.Handler for the list records operation
func NewListRecords(ctx *middleware.Context, handler ListRecordsHandler) *ListRecords {
	return &ListRecords{Context: ctx, Handler: handler}
}

/*ListRecords swagger:route GET /records records listRecords

List records

*/
type ListRecords struct {
	Context *middleware.Context
	Handler ListRecordsHandler
}

func (o *ListRecords) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewListRecordsParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}