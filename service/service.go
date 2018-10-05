package service

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

// ContextKey is a type used for storing values in a context
type ContextKey int

const (
	// ContextKeyUserID is the key used to store user ID in a context
	ContextKeyUserID ContextKey = iota
)

// Opts used to configure the FlameDB service
type Opts struct {
	API            *operations.FlamedbAPI
	Flame          database.Flame
	AccessKeyStore database.AccessKeyStore
	Key            *rsa.PublicKey
}

// Service implements handlers for the FlameDB service
type Service struct {
	api        *operations.FlamedbAPI
	flame      database.Flame
	accessKeys database.AccessKeyStore
	key        *rsa.PublicKey
}

// New returns a new Service instance
func New(opts Opts) *Service {

	svc := Service{
		api:        opts.API,
		flame:      opts.Flame,
		accessKeys: opts.AccessKeyStore,
		key:        opts.Key,
	}

	svc.api.JSONConsumer = runtime.JSONConsumer()
	svc.api.JSONProducer = runtime.JSONProducer()

	svc.api.RecordsCreateRecordHandler = records.CreateRecordHandlerFunc(svc.createRecord)
	svc.api.RecordsDeleteRecordHandler = records.DeleteRecordHandlerFunc(svc.deleteRecord)
	svc.api.RecordsGetRecordHandler = records.GetRecordHandlerFunc(svc.getRecord)
	svc.api.RecordsFindRecordHandler = records.FindRecordHandlerFunc(svc.findRecord)
	svc.api.RecordsListRecordsHandler = records.ListRecordsHandlerFunc(svc.listRecords)
	svc.api.RecordsUpdateRecordHandler = records.UpdateRecordHandlerFunc(svc.updateRecord)

	svc.api.FlamedbAuthAuth = svc.tokenAuth
	svc.api.BasicAuthAuth = svc.basicAuth
	svc.api.ServerShutdown = svc.shutdown

	return &svc
}

func (svc *Service) shutdown() {
	log.Info("Service shutdown")
}

func (svc *Service) tokenAuth(tokenStr string) (*models.Principal, error) {

	// For some reason basic authorization is being processed here in addition
	// to token based auth. For now, differentiate based on the presence of a
	// "Basic " prefix on the token value.

	if strings.HasPrefix(tokenStr, "Basic ") {
		authStr := strings.SplitN(tokenStr, " ", 2)
		payload, err := base64.StdEncoding.DecodeString(authStr[1])
		if err != nil {
			return nil, err
		}
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || pair[0] == "" || pair[1] == "" {
			return nil, errors.New("Invalid basic auth syntax")
		}
		return svc.basicAuth(pair[0], pair[1])
	}

	// If we reach here, then process authentication using JWTs
	if svc.key == nil {
		return nil, errors.New("No token verification key is configured")
	}
	principal, err := parseJWT(svc.key, tokenStr)
	log.Infof("JWT principal: %+v", principal)
	return principal, err
}

func (svc *Service) basicAuth(keyID, keySecret string) (*models.Principal, error) {

	log.Infof("basic auth: %s", keyID)

	if svc.accessKeys == nil {
		return nil, errors.New("Access key store is not configured")
	}

	accessKey, err := svc.accessKeys.Get(keyID)
	if err != nil {
		return nil, err
	}
	log.Infof("Access key: %+v", accessKey)

	if accessKey.Compare(keySecret) != nil {
		return nil, errors.New("Incorrect secret")
	}

	perm := database.AccessKeyPermission(accessKey.Permission)
	isService := perm == database.ServiceRead || perm == database.ServiceReadWrite

	return &models.Principal{
		UserID:      accessKey.RefID,
		Permissions: accessKey.Permission,
		IsService:   isService,
	}, nil
}

func (svc *Service) createRecord(params records.CreateRecordParams, principal *models.Principal) middleware.Responder {

	input := *params.Body

	userID := principal.UserID
	if principal.IsService {
		if params.XUserID == nil || *params.XUserID == "" {
			return records.NewCreateRecordBadRequest().
				WithPayload(&models.BadRequest{
					ErrorType: "BadRequest",
					Message:   "Service must specify X-User-ID",
				})
		}
		userID = *params.XUserID
		log.WithField("user_id", userID).Info("Principal updated")
	}

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
		CreatedBy:  userID,
		UpdatedBy:  userID,
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

	userID := principal.UserID
	if principal.IsService {
		if params.XUserID == nil || *params.XUserID == "" {
			return records.NewDeleteRecordBadRequest().
				WithPayload(&models.BadRequest{
					ErrorType: "BadRequest",
					Message:   "Service must specify X-User-ID",
				})
		}
		userID = *params.XUserID
		log.WithField("user_id", userID).Info("Principal updated")
	}

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

	log.Printf("Record %s deleted by %s", record.ID, userID)
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

	return records.NewGetRecordOK().
		WithPayload(&models.RecordOutput{
			ID:         apiString(record.ID),
			CreatedAt:  apiString(record.CreatedAt.Format(time.RFC3339)),
			CreatedBy:  apiString(record.CreatedBy),
			UpdatedAt:  apiString(record.UpdatedAt.Format(time.RFC3339)),
			UpdatedBy:  apiString(record.UpdatedBy),
			Path:       apiString(record.Path),
			Parent:     apiString(record.Parent),
			Properties: record.MustGetProperties(),
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

	return records.NewFindRecordOK().
		WithPayload(&models.RecordOutput{
			ID:         apiString(record.ID),
			CreatedAt:  apiString(record.CreatedAt.Format(time.RFC3339)),
			CreatedBy:  apiString(record.CreatedBy),
			UpdatedAt:  apiString(record.UpdatedAt.Format(time.RFC3339)),
			UpdatedBy:  apiString(record.UpdatedBy),
			Path:       apiString(record.Path),
			Parent:     apiString(record.Parent),
			Properties: record.MustGetProperties(),
		})
}

func (svc *Service) listRecords(params records.ListRecordsParams, principal *models.Principal) middleware.Responder {

	log.Println("listRecords")

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

	userID := principal.UserID
	if principal.IsService {
		if params.XUserID == nil || *params.XUserID == "" {
			return records.NewUpdateRecordBadRequest().
				WithPayload(&models.BadRequest{
					ErrorType: "BadRequest",
					Message:   "Service must specify X-User-ID",
				})
		}
		userID = *params.XUserID
		log.WithField("user_id", userID).Info("Principal updated")
	}

	propJSON, err := json.Marshal(input.Properties)
	if err != nil {
		return records.NewUpdateRecordBadRequest().
			WithPayload(&models.BadRequest{
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
