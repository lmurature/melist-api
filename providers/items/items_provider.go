package items_provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lmurature/melist-api/domain/apierrors"
	"github.com/lmurature/melist-api/domain/items"
	http_utils "github.com/lmurature/melist-api/utils/http"
	"github.com/mercadolibre/golang-restclient/rest"
	"time"
)

const (
	uriSearchItems = "/sites/MLA/search?q=%s"
)

var (
	itemsRestClient = rest.RequestBuilder{
		BaseURL:        http_utils.BaseUrlMeli,
		Timeout:        5 * time.Second,
		DisableTimeout: false,
	}
)

func SearchItemsByQuery(query string) (*items.ItemSearchResponse, apierrors.ApiError) {
	uri := fmt.Sprintf(uriSearchItems, query)
	fmt.Println(uri)
	response := itemsRestClient.Get(uri)

	if response == nil || response.Response == nil {
		err := errors.New("invalid restclient response")
		msg := "error searching items"
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	if response.StatusCode > 299 {
		apiErr, err := apierrors.NewApiErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error when trying to unmarshal items search response", err)
		}
		return nil, apiErr
	}

	var itemsResult items.ItemSearchResponse
	if err := json.Unmarshal(response.Bytes(), &itemsResult); err != nil {
		msg := "error when trying to unmarshal items search response"
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	return &itemsResult, nil
}
