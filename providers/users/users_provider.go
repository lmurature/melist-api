package users_provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/lmurature/melist-api/domain/apierrors"
	"github.com/lmurature/melist-api/domain/users"
	http_utils "github.com/lmurature/melist-api/utils/http"
	"github.com/mercadolibre/golang-restclient/rest"
)

const (
	getUserUri = "/users/%d?caller.id=%d"
	BEARER     = "Bearer %s"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL:        http_utils.BaseUrlMeli,
		Timeout:        5 * time.Second,
		DisableTimeout: false,
	}
)

func GetUserInformation(userId int64, accessToken string) (*users.User, apierrors.ApiError) {
	usersRestClient.Headers = make(http.Header)
	usersRestClient.Headers.Add("Authorization", fmt.Sprintf(BEARER, accessToken))
	defer usersRestClient.Headers.Del("Authorization")

	uri := fmt.Sprintf(getUserUri, userId, userId)

	response := usersRestClient.Get(uri)

	if response == nil || response.Response == nil {
		msg := fmt.Sprintf("invalid restclient response while trying to get information for user %d", userId)
		err := errors.New("invalid restclient response")
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	if response.StatusCode > 299 {
		apiErr, err := apierrors.NewApiErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error when trying to unmarshal get user response", err)
		}
		return nil, apiErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		msg := fmt.Sprintf("error while trying to unmarshal data from user %d", userId)
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	return &user, nil
}
