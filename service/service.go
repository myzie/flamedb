package service

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/myzie/flamedb/restapi/operations"
	"github.com/myzie/flamedb/restapi/operations/records"
)

// Opts used to configure the FlameDB service
type Opts struct {
	API *operations.FlamedbAPI
}

// Service implements handlers for the FlameDB service
type Service struct {
	api *operations.FlamedbAPI
}

// New returns a new Service instance
func New(opts Opts) *Service {

	svc := Service{
		api: opts.API,
	}

	svc.api.JSONConsumer = runtime.JSONConsumer()
	svc.api.JSONProducer = runtime.JSONProducer()

	// auth := jwtAuth{
	// 	Key: nil,
	// 	AuthFunction: func(claims jwt.Claims) (interface{}, error) {
	// 		sub, ok := claims["sub"]
	// 		if !ok {
	// 			return nil, errors.New("Sub not provided")
	// 		}
	// 		return sub, nil
	// 	},
	// }

	svc.api.RecordsCreateRecordHandler = records.CreateRecordHandlerFunc(svc.CreateRecord)
	svc.api.RecordsDeleteRecordHandler = records.DeleteRecordHandlerFunc(svc.DeleteRecord)
	svc.api.RecordsGetRecordHandler = records.GetRecordHandlerFunc(svc.GetRecord)
	svc.api.RecordsListRecordsHandler = records.ListRecordsHandlerFunc(svc.ListRecords)
	svc.api.RecordsUpdateRecordHandler = records.UpdateRecordHandlerFunc(svc.UpdateRecord)

	svc.api.ServerShutdown = svc.shutdown

	return &svc
}

func (svc *Service) shutdown() {

}

// CreateRecord ...
func (svc *Service) CreateRecord(params records.CreateRecordParams, principal interface{}) middleware.Responder {
	return records.NewCreateRecordOK()
}

// DeleteRecord ...
func (svc *Service) DeleteRecord(params records.DeleteRecordParams, principal interface{}) middleware.Responder {
	return nil
}

// GetRecord ...
func (svc *Service) GetRecord(params records.GetRecordParams, principal interface{}) middleware.Responder {
	return records.NewGetRecordOK()
}

// ListRecords ...
func (svc *Service) ListRecords(params records.ListRecordsParams, principal interface{}) middleware.Responder {
	return records.NewListRecordsOK()
}

// UpdateRecord ...
func (svc *Service) UpdateRecord(params records.UpdateRecordParams, principal interface{}) middleware.Responder {
	return records.NewUpdateRecordOK()
}
