package items_provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lmurature/golang-restclient/rest"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/items"
	http_utils "github.com/lmurature/melist-api/src/api/utils/http"
	"strconv"
	"strings"
	"time"
)

const (
	uriSearchItems        = "/sites/MLA/search?q=%s&offset=%d"
	uriGetItem            = "/items/%s"
	uriGetItemDescription = "/items/%s/description"
	uriGetItemReviews     = "/reviews/item/%s?catalog_product_id=%s&limit=200&order=desc&order_criteria=dateCreated"
	uriGetCategoryTrends  = "/trends/MLA/%s"
	uriGetCategory        = "/categories/%s"
)

var (
	itemsRestClient = rest.RequestBuilder{
		BaseURL:        http_utils.BaseUrlMeli,
		Timeout:        15 * time.Second,
		DisableTimeout: false,
	}

	reviewsRestClient = rest.RequestBuilder{
		BaseURL:        http_utils.BaseUrlMeli,
		Timeout:        5 * time.Second,
		DisableCache:   true,
		DisableTimeout: false,
	}

	vipRestClient = rest.RequestBuilder{
		Timeout:        5 * time.Second,
		DisableCache:   true,
		DisableTimeout: false,
	}

	categoryRestClient = rest.RequestBuilder{
		BaseURL:        http_utils.BaseUrlMeli,
		Timeout:        15 * time.Second,
		DisableTimeout: false,
		DisableCache:   false,
	}
)

func SearchItemsByQuery(query string, offset int) (*items.ItemSearchResponse, apierrors.ApiError) {
	uri := fmt.Sprintf(uriSearchItems, query, offset)
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

func GetItemReviews(itemId string, catalogProductId string) (*items.ItemReviewsResponse, apierrors.ApiError) {
	uri := fmt.Sprintf(uriGetItemReviews, itemId, catalogProductId)
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

func GetCategoryTrends(categoryId string) (*items.CategoryTrends, apierrors.ApiError) {
	uri := fmt.Sprintf(uriGetCategoryTrends, categoryId)
	response := itemsRestClient.Get(uri)

	if response == nil || response.Response == nil {
		err := errors.New("invalid restclient response")
		msg := fmt.Sprintf("invalid restclient response getting category %s trends", categoryId)
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	if response.StatusCode > 299 {
		apiErr, err := apierrors.NewApiErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error when trying to unmarshal apierror category trends response", err)
		}
		return nil, apiErr
	}

	var result items.CategoryTrends
	if err := json.Unmarshal(response.Bytes(), &result); err != nil {
		msg := fmt.Sprintf("error trying to unmarshal response into category %s trends structure", categoryId)
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	return &result, nil
}

func GetRealQuantity(permalink string) (*int64, apierrors.ApiError) {
	response := vipRestClient.Get(permalink)

	if response == nil || response.Response == nil {
		return nil, apierrors.NewInternalServerApiError("nil resp", errors.New("nil resp"))
	}

	if response.StatusCode > 299 {
		apiErr, err := apierrors.NewApiErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error when trying to unmarshal apierror category trends response", err)
		}
		return nil, apiErr
	}

	var result = response.String()

	values := strings.Split(result, "\"availableStock\":")
	if len(values) > 1 {
		quantityString := strings.Split(values[1], ",")
		quantity, err := strconv.ParseInt(quantityString[0], 10, 32)
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error getting int from string", err)
		}
		return &quantity, nil
	}

	return nil, nil
}

func GetCategory(categoryId string) (*items.Category, apierrors.ApiError) {
	uri := fmt.Sprintf(uriGetCategory, categoryId)
	response := categoryRestClient.Get(uri)

	if response == nil || response.Response == nil {
		err := errors.New("invalid restclient response")
		msg := fmt.Sprintf("invalid restclient response getting category %s", categoryId)
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	if response.StatusCode > 299 {
		apiErr, err := apierrors.NewApiErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, apierrors.NewInternalServerApiError("error when trying to unmarshal apierror category", err)
		}
		return nil, apiErr
	}

	var cat items.Category
	if err := json.Unmarshal(response.Bytes(), &cat); err != nil {
		msg := fmt.Sprintf("error trying to unmarshal response into category %s structure", categoryId)
		return nil, apierrors.NewInternalServerApiError(msg, err)
	}

	return &cat, nil
}
