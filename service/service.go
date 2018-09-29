package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/myzie/flamedb/database"
	"github.com/myzie/flamedb/models"
	"github.com/myzie/flamedb/restapi/operations"
	"github.com/myzie/flamedb/restapi/operations/records"
	log "github.com/sirupsen/logrus"
)

// Opts used to configure the FlameDB service
type Opts struct {
	API   *operations.FlamedbAPI
	Flame database.Flame
}

// Service implements handlers for the FlameDB service
type Service struct {
	api   *operations.FlamedbAPI
	flame database.Flame
}

// New returns a new Service instance
func New(opts Opts) *Service {

	svc := Service{
		api:   opts.API,
		flame: opts.Flame,
	}

	svc.api.JSONConsumer = runtime.JSONConsumer()
	svc.api.JSONProducer = runtime.JSONProducer()

	svc.api.RecordsCreateRecordHandler = records.CreateRecordHandlerFunc(svc.createRecord)
	svc.api.RecordsDeleteRecordHandler = records.DeleteRecordHandlerFunc(svc.deleteRecord)
	svc.api.RecordsGetRecordHandler = records.GetRecordHandlerFunc(svc.getRecord)
	svc.api.RecordsFindRecordHandler = records.FindRecordHandlerFunc(svc.findRecord)
	svc.api.RecordsListRecordsHandler = records.ListRecordsHandlerFunc(svc.listRecords)
	svc.api.RecordsUpdateRecordHandler = records.UpdateRecordHandlerFunc(svc.updateRecord)

	svc.api.FlamedbAuthAuth = svc.authenticate
	svc.api.ServerShutdown = svc.shutdown

	return &svc
}

func (svc *Service) shutdown() {
	log.Info("Service shutdown")
}

func (svc *Service) authenticate(token string) (*models.Principal, error) {
	log.Infof("authenticate: %s", token)

	records, err := svc.flame.List(database.Query{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Infof("Records: %+v\n", records)

	return nil, nil
}

func (svc *Service) createRecord(params records.CreateRecordParams, principal *models.Principal) middleware.Responder {

	input := *params.Body

	if _, err := svc.flame.Get(database.Key{Path: *input.Path}); err == nil {
		return records.NewCreateRecordBadRequest().
			WithPayload(&models.BadRequest{
				ErrorType: "BadRequest",
				Message:   "Record already exists at that path",
			})
	}

	propJSON, err := json.Marshal(input.Properties)
	if err != nil {
		return records.NewCreateRecordBadRequest().
			WithPayload(&models.BadRequest{
				ErrorType: "ValidationError",
				Message:   "Invalid properties JSON",
			})
	}

	record := database.Record{
		Path:       *input.Path,
		Properties: postgres.Jsonb{RawMessage: json.RawMessage(propJSON)},
	}
	if err = svc.flame.Save(&record); err != nil {
		log.WithError(err).Info("Save error")
		return records.NewCreateRecordInternalServerError().
			WithPayload(&models.InternalServerError{
				ErrorType: "InternalServerError",
				Message:   "Failed to save record",
			})
	}

	return records.NewCreateRecordOK().
		WithPayload(&models.RecordOutput{
			ID:         apiString(record.ID),
			CreatedAt:  apiString(record.CreatedAt.Format(time.RFC3339)),
			CreatedBy:  apiString(record.CreatedBy),
			UpdatedAt:  apiString(record.UpdatedAt.Format(time.RFC3339)),
			UpdatedBy:  apiString(record.UpdatedBy),
			Path:       apiString(record.Path),
			Parent:     apiString(record.Parent),
			Properties: input.Properties,
		})
}

func (svc *Service) deleteRecord(params records.DeleteRecordParams, principal *models.Principal) middleware.Responder {

	record, err := svc.flame.Get(database.Key{ID: params.RecordID})
	if err != nil {
		log.WithError(err).Info("Get error")
		return records.NewDeleteRecordNotFound().
			WithPayload(&models.NotFoundError{
				ErrorType: "NotFound",
				Message:   "Record not found",
			})
	}

	if err := svc.flame.Delete(record); err != nil {
		log.WithError(err).Info("Delete error")
		return records.NewDeleteRecordInternalServerError().
			WithPayload(&models.InternalServerError{
				ErrorType: "InternalServerError",
				Message:   "Failed to delete record",
			})
	}
	return records.NewDeleteRecordOK()
}

func (svc *Service) getRecord(params records.GetRecordParams, principal *models.Principal) middleware.Responder {

	record, err := svc.flame.Get(database.Key{ID: params.RecordID})
	if err != nil {
		log.WithError(err).Info("Get error")
		return records.NewGetRecordNotFound().
			WithPayload(&models.NotFoundError{
				ErrorType: "NotFound",
				Message:   "Record not found",
			})
	}

	properties, err := record.GetProperties()
	if err != nil {
		log.WithError(err).Error("Failed to marshal record JSON")
		return records.NewGetRecordInternalServerError().
			WithPayload(&models.InternalServerError{
				ErrorType: "InternalServerError",
				Message:   "Failed to get record properties",
			})
	}

	return records.NewGetRecordOK().
		WithPayload(&models.RecordOutput{
			ID:         apiString(record.ID),
			CreatedAt:  apiString(record.CreatedAt.Format(time.RFC3339)),
			CreatedBy:  apiString(record.CreatedBy),
			UpdatedAt:  apiString(record.UpdatedAt.Format(time.RFC3339)),
			UpdatedBy:  apiString(record.UpdatedBy),
			Path:       apiString(record.Path),
			Parent:     apiString(record.Parent),
			Properties: properties,
		})
}

func (svc *Service) findRecord(params records.FindRecordParams, principal *models.Principal) middleware.Responder {

	record, err := svc.flame.Get(database.Key{Path: params.Path})
	if err != nil {
		log.WithError(err).Info("Get error")
		return records.NewFindRecordNotFound().
			WithPayload(&models.NotFoundError{
				ErrorType: "NotFound",
				Message:   "Record not found",
			})
	}

	properties, err := record.GetProperties()
	if err != nil {
		log.WithError(err).Error("Failed to marshal record JSON")
		return records.NewFindRecordInternalServerError().
			WithPayload(&models.InternalServerError{
				ErrorType: "InternalServerError",
				Message:   "Failed to get record properties",
			})
	}

	return records.NewFindRecordOK().
		WithPayload(&models.RecordOutput{
			ID:         apiString(record.ID),
			CreatedAt:  apiString(record.CreatedAt.Format(time.RFC3339)),
			CreatedBy:  apiString(record.CreatedBy),
			UpdatedAt:  apiString(record.UpdatedAt.Format(time.RFC3339)),
			UpdatedBy:  apiString(record.UpdatedBy),
			Path:       apiString(record.Path),
			Parent:     apiString(record.Parent),
			Properties: properties,
		})
}

func (svc *Service) listRecords(params records.ListRecordsParams, principal *models.Principal) middleware.Responder {

	query := database.Query{
		Offset:              getIntDefault(params.Offset, 0),
		Limit:               getIntDefault(params.Limit, 100),
		Parent:              getStrDefault(params.Parent, ""),
		Prefix:              getStrDefault(params.Prefix, "/"),
		OrderBy:             getStrDefault(params.OrderBy, ""),
		OrderByDesc:         getBoolDefault(params.OrderByDesc, false),
		OrderByProperty:     getStrDefault(params.OrderByProperty, ""),
		OrderByPropertyDesc: getBoolDefault(params.OrderByPropertyDesc, false),
	}
	results, err := svc.flame.List(query)
	if err != nil {
		return records.NewListRecordsInternalServerError()
	}

	items := make([]*models.RecordOutput, len(results))
	for i, r := range results {

		createdAt := r.CreatedAt.Format(time.RFC3339)
		updatedAt := r.UpdatedAt.Format(time.RFC3339)

		var props map[string]interface{}
		if err := json.Unmarshal(r.Properties.RawMessage, &props); err != nil {
			return records.NewListRecordsInternalServerError()
		}

		items[i] = &models.RecordOutput{
			ID:         apiString(r.ID),
			Parent:     apiString(r.Parent),
			Path:       apiString(r.Path),
			CreatedAt:  apiString(createdAt),
			CreatedBy:  apiString(r.CreatedBy),
			UpdatedAt:  apiString(updatedAt),
			UpdatedBy:  apiString(r.UpdatedBy),
			Properties: props,
		}
	}

	return records.NewListRecordsOK().
		WithPayload(&models.QueryResult{Records: items})
}

func (svc *Service) updateRecord(params records.UpdateRecordParams, principal *models.Principal) middleware.Responder {

	input := *params.Record

	propJSON, err := json.Marshal(input.Properties)
	if err != nil {
		return records.NewUpdateRecordBadRequest().
			WithPayload(&models.ValidationError{
				ErrorType: "ValidationError",
				Message:   "Invalid properties JSON",
			})
	}

	record, err := svc.flame.Get(database.Key{ID: params.RecordID})
	if err != nil {
		return records.NewUpdateRecordNotFound().
			WithPayload(&models.NotFoundError{
				ErrorType: "NotFound",
				Message:   fmt.Sprintf("Record with ID %s was not found", params.RecordID),
			})
	}

	record.Properties = postgres.Jsonb{RawMessage: json.RawMessage(propJSON)}
	if err = svc.flame.Save(record); err != nil {
		return records.NewUpdateRecordInternalServerError().
			WithPayload(&models.InternalServerError{
				ErrorType: "InternalServerError",
				Message:   "Failed to update record",
			})
	}

	return records.NewUpdateRecordOK().WithPayload(&models.RecordOutput{
		ID:         apiString(record.ID),
		CreatedAt:  apiString(record.CreatedAt.Format(time.RFC3339)),
		CreatedBy:  apiString(record.CreatedBy),
		UpdatedAt:  apiString(record.UpdatedAt.Format(time.RFC3339)),
		UpdatedBy:  apiString(record.UpdatedBy),
		Path:       apiString(record.Path),
		Parent:     apiString(record.Parent),
		Properties: input.Properties,
	})
}

func getStrDefault(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

func getIntDefault(i *int64, def int) int {
	if i == nil {
		return def
	}
	return int(*i)
}

func getBoolDefault(b *bool, def bool) bool {
	if b == nil {
		return def
	}
	return *b
}

func apiString(s string) *string {
	return &s
}
