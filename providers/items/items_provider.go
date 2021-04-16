package items_provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lmurature/golang-restclient/rest"
	"github.com/lmurature/melist-api/domain/apierrors"
	"github.com/lmurature/melist-api/domain/items"
	http_utils "github.com/lmurature/melist-api/utils/http"
	"time"
)

const (
	uriSearchItems = "/sites/MLA/search?q=%s"
	uriGetItem     = "/items/%s"
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
	response := itemsRestClient.Get(uri)

	if response == nil || response.Response == nil {
		err := errors.New("invalid restclient response")
		msg := "invalid restclient response searching items"
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	if response.StatusCode > 299 {
		apiErr, err := apierrors.NewApiErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error when trying to unmarshal apierror items search response", err)
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

func GetItemById(itemId string) (*items.Item, apierrors.ApiError) {
	uri := fmt.Sprintf(uriGetItem, itemId)
	response := itemsRestClient.Get(uri)

	if response == nil || response.Response == nil {
		err := errors.New("invalid restclient response")
		msg := fmt.Sprintf("invalid restclient response getting item %s", itemId)
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	if response.StatusCode > 299 {
		apiErr, err := apierrors.NewApiErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error when trying to unmarshal apierror item response", err)
		}
		return nil, apiErr
	}

	var item items.Item
	if err := json.Unmarshal(response.Bytes(), &item); err != nil {
		msg := fmt.Sprintf("error trying to unmarshal response into item %s structure", itemId)
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	return &item, nil
}
