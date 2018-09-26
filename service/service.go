package service

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/myzie/flamedb/models"
	"github.com/myzie/flamedb/restapi/operations"
	"github.com/myzie/flamedb/restapi/operations/records"
	log "github.com/sirupsen/logrus"
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

	svc.api.RecordsCreateRecordHandler = records.CreateRecordHandlerFunc(svc.createRecord)
	svc.api.RecordsDeleteRecordHandler = records.DeleteRecordHandlerFunc(svc.deleteRecord)
	svc.api.RecordsGetRecordHandler = records.GetRecordHandlerFunc(svc.getRecord)
	svc.api.RecordsListRecordsHandler = records.ListRecordsHandlerFunc(svc.listRecords)
	svc.api.RecordsUpdateRecordHandler = records.UpdateRecordHandlerFunc(svc.updateRecord)

	svc.api.FlamedbAuthAuth = svc.authenticate
	svc.api.ServerShutdown = svc.shutdown

	log.Info("Service created")
	return &svc
}

func (svc *Service) shutdown() {
	log.Info("Service shutdown")
}

func (svc *Service) authenticate(token string) (*models.Principal, error) {
	log.Infof("authenticate: %s", token)
	return nil, nil
}

// CreateRecord ...
func (svc *Service) createRecord(params records.CreateRecordParams, principal *models.Principal) middleware.Responder {
	log.Info("createRecord")
	return records.NewCreateRecordOK()
}

// DeleteRecord ...
func (svc *Service) deleteRecord(params records.DeleteRecordParams, principal *models.Principal) middleware.Responder {
	log.Info("deleteRecord")
	return nil
}

// GetRecord ...
func (svc *Service) getRecord(params records.GetRecordParams, principal *models.Principal) middleware.Responder {
	log.Info("getRecord")
	return records.NewGetRecordOK()
}

// ListRecords ...
func (svc *Service) listRecords(params records.ListRecordsParams, principal *models.Principal) middleware.Responder {
	log.Info("listRecords")
	return records.NewListRecordsOK()
}

// UpdateRecord ...
func (svc *Service) updateRecord(params records.UpdateRecordParams, principal *models.Principal) middleware.Responder {
	log.Info("updateRecord")
	return records.NewUpdateRecordOK()
}
