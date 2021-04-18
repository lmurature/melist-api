package users_provider

import (
	"github.com/lmurature/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestGetUserInformationInvalidResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: -1,
	})

	user, err := GetUserInformation(1)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "invalid restclient response while trying to get information for user 1", err.Message())
}

func TestGetUserInformationErrorInvalidApiError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusForbidden,
		RespBody:     `{---`,
	})

	user, err := GetUserInformation(1)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error when trying to unmarshal get user response error to ApiError", err.Message())
}

func TestGetUserInformationErrorFromBytes(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"error": "user not found", "message": "this user does not exist", "status": 404}`,
	})

	user, err := GetUserInformation(1)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
	assert.EqualValues(t, "this user does not exist", err.Message())
}

func TestGetUserInformationErrorUnmarshalIntoUser(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{---`,
	})

	user, err := GetUserInformation(1)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error while trying to unmarshal data from user 1", err.Message())
}

func TestGetUserInformationNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 1, "nickname": "pepe"}`,
	})

	user, err := GetUserInformation(1)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "pepe", user.Nickname)
}

func TestGetUserInformationMeInvalidResponse(t *testing.T) {
	rest.FlushMockups()
	headers := make(http.Header)
	headers.Add("Authorization", "Bearer a1b2c3d4e5")
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/me",
		HTTPMethod:   http.MethodGet,
		ReqHeaders:   headers,
		RespHTTPCode: -1,
	})

	user, err := GetUserInformationMe("a1b2c3d4e5")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "invalid restclient response while trying to get information for my user", err.Message())
}

func TestGetUserInformationMeErrorInvalidApiError(t *testing.T) {
	rest.FlushMockups()
	headers := make(http.Header)
	headers.Add("Authorization", "Bearer a1b2c3d4e5")
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/me",
		HTTPMethod:   http.MethodGet,
		ReqHeaders:   headers,
		RespHTTPCode: http.StatusForbidden,
		RespBody:     `{---`,
	})

	user, err := GetUserInformationMe("a1b2c3d4e5")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error when trying to unmarshal get user response error to ApiError", err.Message())
}

func TestGetUserInformationMeErrorFromBytes(t *testing.T) {
	rest.FlushMockups()
	headers := make(http.Header)
	headers.Add("Authorization", "Bearer a1b2c3d4e5")
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/me",
		HTTPMethod:   http.MethodGet,
		ReqHeaders:   headers,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"error": "user not found", "message": "this user does not exist", "status": 404}`,
	})

	user, err := GetUserInformationMe("a1b2c3d4e5")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
	assert.EqualValues(t, "this user does not exist", err.Message())
}

func TestGetUserInformationMeErrorUnmarshalIntoUser(t *testing.T) {
	rest.FlushMockups()
	headers := make(http.Header)
	headers.Add("Authorization", "Bearer a1b2c3d4e5")
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/me",
		HTTPMethod:   http.MethodGet,
		ReqHeaders:   headers,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{---`,
	})

	user, err := GetUserInformationMe("a1b2c3d4e5")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error while trying to unmarshal data from my user", err.Message())
}

func TestGetUserInformationMeNoError(t *testing.T) {
	rest.FlushMockups()
	headers := make(http.Header)
	headers.Add("Authorization", "Bearer a1b2c3d4e5")
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/me",
		HTTPMethod:   http.MethodGet,
		ReqHeaders:   headers,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 1, "nickname": "pepe"}`,
	})

	user, err := GetUserInformationMe("a1b2c3d4e5")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "pepe", user.Nickname)
}

