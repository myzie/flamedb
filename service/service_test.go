package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/myzie/flamedb/database"

	"github.com/go-openapi/loads"
	"github.com/golang/mock/gomock"
	"github.com/myzie/flamedb/database/mock_database"
	"github.com/myzie/flamedb/restapi"
	"github.com/myzie/flamedb/restapi/operations"
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
