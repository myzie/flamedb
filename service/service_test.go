package service

import (
	"errors"
	"testing"
	"time"

	"github.com/myzie/flamedb/models"
	"github.com/myzie/flamedb/restapi/operations/records"

	"github.com/go-openapi/loads"
	"github.com/golang/mock/gomock"
	"github.com/myzie/flamedb/database"
	"github.com/myzie/flamedb/database/mock_database"
	"github.com/myzie/flamedb/restapi"
	"github.com/myzie/flamedb/restapi/operations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestAPI() *operations.FlamedbAPI {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		panic(err)
	}
	return operations.NewFlamedbAPI(swaggerSpec)
}

func TestServiceInit(t *testing.T) {
	svc := New(Opts{API: getTestAPI()})
	require.NotNil(t, svc)
}

func TestAccessKeyOK(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	accessKey, secret, err := database.NewAccessKey(
		"test-key",
		"ref-id",
		database.ReadWrite,
	)
	require.Nil(t, err)

	mockKeyStore := mock_database.NewMockAccessKeyStore(mockCtrl)
	mockKeyStore.EXPECT().Get("key-id").Return(accessKey, nil).Times(2)

	svc := New(Opts{API: getTestAPI(), AccessKeyStore: mockKeyStore})
	require.NotNil(t, svc)

	p, err := svc.basicAuth("key-id", secret)
	require.Nil(t, err)
	require.NotNil(t, p)
	assert.False(t, p.IsService)
	assert.Equal(t, p.Permissions, string(database.ReadWrite))
	assert.Equal(t, p.UserID, "ref-id")

	p, err = svc.basicAuth("key-id", "wrong")
	require.NotNil(t, err)
	require.Nil(t, p)
	assert.Equal(t, err.Error(), "Incorrect secret")
}

func TestAccessKeyNotFound(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockKeyStore := mock_database.NewMockAccessKeyStore(mockCtrl)
	mockKeyStore.EXPECT().Get("key-id").Return(nil, errors.New("bad key"))

	svc := New(Opts{API: getTestAPI(), AccessKeyStore: mockKeyStore})
	require.NotNil(t, svc)

	p, err := svc.basicAuth("key-id", "secret")
	require.NotNil(t, err)
	require.Nil(t, p)
	assert.Equal(t, err.Error(), "bad key")
}

func TestCreateRecord(t *testing.T) {

	assert := assert.New(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	now := time.Now()
	nowStr := now.Format(time.RFC3339)

	flame := mock_database.NewMockFlame(mockCtrl)
	flame.EXPECT().Get(database.Key{Path: "/items/1"}).Return(nil, errors.New(""))
	flame.EXPECT().Save(gomock.Any()).Return(nil).Do(func(r *database.Record) {
		assert.Equal("", r.ID)
		assert.True(r.CreatedAt.IsZero())
		assert.True(r.UpdatedAt.IsZero())
		assert.Equal("/items/1", r.Path)
		assert.Equal("user-1", r.CreatedBy)
		assert.Equal("user-1", r.UpdatedBy)
		assert.Equal(map[string]interface{}{"foo": "bar"}, r.MustGetProperties())

		r.ID = "new-id"
		r.CreatedAt = now
		r.UpdatedAt = now
		r.Parent = "/items/"
	})

	svc := New(Opts{API: getTestAPI(), Flame: flame})
	require.NotNil(t, svc)

	params := records.CreateRecordParams{
		Body: &models.RecordInput{
			Path:       apiString("/items/1"),
			Properties: map[string]interface{}{"foo": "bar"},
		},
	}
	principal := &models.Principal{UserID: "user-1"}

	resp, ok := svc.createRecord(params, principal).(*records.CreateRecordOK)
	require.True(t, ok)
	out := resp.Payload
	require.NotNil(t, out)

	assert.Equal("new-id", *out.ID)
	assert.Equal(nowStr, *out.CreatedAt)
	assert.Equal(nowStr, *out.UpdatedAt)
	assert.Equal("user-1", *out.CreatedBy)
	assert.Equal("user-1", *out.UpdatedBy)
	assert.Equal("/items/", *out.Parent)
	assert.Equal("/items/1", *out.Path)
	assert.Equal(map[string]interface{}{"foo": "bar"}, out.Properties)
}

func TestDeleteRecord(t *testing.T) {

	assert := assert.New(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	record := &database.Record{
		Path: "/path/to/item",
	}

	flame := mock_database.NewMockFlame(mockCtrl)
	flame.EXPECT().Get(database.Key{ID: "id-1"}).Return(record, nil)
	flame.EXPECT().Delete(gomock.Any()).Return(nil).Do(func(r *database.Record) {
		assert.Equal("", r.ID)
		assert.Equal("/path/to/item", r.Path)
	})

	svc := New(Opts{API: getTestAPI(), Flame: flame})
	require.NotNil(t, svc)

	params := records.DeleteRecordParams{RecordID: "id-1"}
	principal := &models.Principal{UserID: "user-1"}

	_, ok := svc.deleteRecord(params, principal).(*records.DeleteRecordOK)
	require.True(t, ok)
}
