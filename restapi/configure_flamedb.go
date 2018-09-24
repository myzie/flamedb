// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/myzie/flamedb/restapi/operations"
	"github.com/myzie/flamedb/restapi/operations/records"
)

//go:generate swagger generate server --target .. --name flamedb --spec ../swagger.yaml

func configureFlags(api *operations.FlamedbAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.FlamedbAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "api_key" header is set
	api.FlamedbAuthAuth = func(token string) (interface{}, error) {
		return nil, errors.NotImplemented("api key auth (flamedb_auth) api_key from header param [api_key] has not yet been implemented")
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	api.RecordsCreateRecordHandler = records.CreateRecordHandlerFunc(func(params records.CreateRecordParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation records.CreateRecord has not yet been implemented")
	})
	api.RecordsDeleteRecordHandler = records.DeleteRecordHandlerFunc(func(params records.DeleteRecordParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation records.DeleteRecord has not yet been implemented")
	})
	api.RecordsGetRecordHandler = records.GetRecordHandlerFunc(func(params records.GetRecordParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation records.GetRecord has not yet been implemented")
	})
	api.RecordsListRecordsHandler = records.ListRecordsHandlerFunc(func(params records.ListRecordsParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation records.ListRecords has not yet been implemented")
	})
	api.RecordsUpdateRecordHandler = records.UpdateRecordHandlerFunc(func(params records.UpdateRecordParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation records.UpdateRecord has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
