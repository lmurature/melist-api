package items_provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lmurature/golang-restclient/rest"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/items"
	http_utils "github.com/lmurature/melist-api/src/api/utils/http"
	"time"
)

const (
	uriSearchItems        = "/sites/MLA/search?q=%s"
	uriGetItem            = "/items/%s"
	uriGetItemDescription = "/items/%s/description"
	uriGetItemReviews     = "/reviews/item/%s"
)

var (
	itemsRestClient = rest.RequestBuilder{
		BaseURL:        http_utils.BaseUrlMeli,
		Timeout:        5 * time.Second,
		DisableTimeout: false,
	}

	reviewsRestClient = rest.RequestBuilder{
		BaseURL:        http_utils.BaseUrlMeli,
		Timeout:        5 * time.Second,
		DisableCache:   true,
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

func GetItemDescription(itemId string) (*items.ItemDescription, apierrors.ApiError) {
	uri := fmt.Sprintf(uriGetItemDescription, itemId)
	response := itemsRestClient.Get(uri)

	if response == nil || response.Response == nil {
		err := errors.New("invalid restclient response")
		msg := fmt.Sprintf("invalid restclient response getting item %s description", itemId)
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	if response.StatusCode > 299 {
		apiErr, err := apierrors.NewApiErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error when trying to unmarshal apierror item description response", err)
		}
		return nil, apiErr
	}

	var description items.ItemDescription
	if err := json.Unmarshal(response.Bytes(), &description); err != nil {
		msg := fmt.Sprintf("error trying to unmarshal response into item description %s structure", itemId)
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	return &description, nil
}

func GetItemReviews(itemId string) (*items.ItemReviewsResponse, apierrors.ApiError) {
	uri := fmt.Sprintf(uriGetItemReviews, itemId)
	response := reviewsRestClient.Get(uri)

	if response == nil || response.Response == nil {
		err := errors.New("invalid restclient response")
		msg := fmt.Sprintf("invalid restclient response getting item %s reviews", itemId)
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	if response.StatusCode > 299 {
		apiErr, err := apierrors.NewApiErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error when trying to unmarshal apierror item review response", err)
		}
		return nil, apiErr
	}

	var result items.ItemReviewsResponse
	if err := json.Unmarshal(response.Bytes(), &result); err != nil {
		msg := fmt.Sprintf("error trying to unmarshal response into item review %s structure", itemId)
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	return &result, nil
}
