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
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"error": "not_found", "message": "user not found", "status": 404}`,
	})

	s := &usersService{}

	user, err := s.GetMeliUser(1)

	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.EqualValues(t, "user not found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.Status())
}

func TestGetUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.mercadolibre.com/users/1",
		HTTPMethod:   http.MethodGet,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 1, "nickname": "pepe"}`,
	})

	s := &usersService{}

	user, err := s.GetMeliUser(1)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "pepe", user.Nickname)
}
