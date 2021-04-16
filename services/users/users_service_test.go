package users_service

import (
	"net/http"
	"os"
	"testing"

	"github.com/lmurature/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestGetUserError(t *testing.T) {
	rest.FlushMockups()
	headers := make(http.Header)
	headers.Add("Authorization", "Bearer a1b2c3d4e5")
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/1",
		HTTPMethod:   http.MethodGet,
		ReqHeaders:   headers,
		RespHTTPCode: http.StatusForbidden,
		RespBody:     `{"error": "forbidden", "message": "you are not authorized to access this resource", "status": 403}`,
	})

	s := &usersService{}

	user, err := s.GetUser(1, "a1b2c3d4e5")

	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.EqualValues(t, "you are not authorized to access this resource", err.Message())
	assert.EqualValues(t, 403, err.Status())
}

func TestGetUserNoError(t *testing.T) {
	rest.FlushMockups()
	headers := make(http.Header)
	headers.Add("Authorization", "Bearer a1b2c3d4e5")
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/1",
		HTTPMethod:   http.MethodGet,
		ReqHeaders:   headers,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 1, "nickname": "pepe"}`,
	})

	s := &usersService{}

	user, err := s.GetUser(1, "a1b2c3d4e5")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "pepe", user.Nickname)
}
