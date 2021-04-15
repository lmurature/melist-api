package auth_provider

import (
	"encoding/json"
	"errors"
	"github.com/lmurature/melist-api/config"
	"github.com/lmurature/melist-api/domain/apierrors"
	"github.com/lmurature/melist-api/domain/auth"
	http_utils "github.com/lmurature/melist-api/utils/http"
	"github.com/mercadolibre/golang-restclient/rest"
	"time"
)

const (
	uriAuthenticateUser = "/oauth/token"
)

var (
	authenticationRestClient = rest.RequestBuilder{
		BaseURL:        http_utils.BaseUrlMeli,
		Timeout:        5 * time.Second,
		DisableTimeout: false,
	}
)

func CreateUserAccessToken(code string) (*auth.MeliAuthResponse, apierrors.ApiError) {
	requestBody := auth.MeliAuthRequest{
		GrantType:    auth.GrantTypeAuthorizationCode,
		ClientId:     config.AppId,
		ClientSecret: config.SecretKey,
		Code:         code,
		RedirectUri:  config.RedirectUri,
	}

	response := authenticationRestClient.Post(uriAuthenticateUser, requestBody)

	if response == nil || response.Response == nil {
		err := errors.New("invalid restclient response")
		msg := "error authenticating user"
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	if response.StatusCode > 299 {
		apiErr, err := apierrors.NewApiErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error when trying to unmarshal authenticate user response", err)
		}
		return nil, apiErr
	}

	var result auth.MeliAuthResponse
	if err := json.Unmarshal(response.Bytes(), &result); err != nil {
		msg := "error when trying to unmarshal user auth response"
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	return &result, nil
}

func RefreshAccessToken(refreshToken string) (*auth.MeliAuthResponse, apierrors.ApiError) {
	requestBody := auth.MeliAuthRequest{
		GrantType:    auth.GrantTypeRefreshToken,
		ClientId:     config.AppId,
		ClientSecret: config.SecretKey,
		RefreshToken: refreshToken,
	}

	response := authenticationRestClient.Post(uriAuthenticateUser, requestBody)

	if response == nil || response.Response == nil {
		err := errors.New("invalid restclient response")
		msg := "error authenticating user"
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	if response.StatusCode > 299 {
		apiErr, err := apierrors.NewApiErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error when trying to unmarshal authenticate user response", err)
		}
		return nil, apiErr
	}

	var result auth.MeliAuthResponse
	if err := json.Unmarshal(response.Bytes(), &result); err != nil {
		msg := "error when trying to unmarshal user refresh token response"
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	return &result, nil
}
