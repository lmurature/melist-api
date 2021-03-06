package items_service

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"github.com/lmurature/melist-api/src/api/domain/items"
	items_provider "github.com/lmurature/melist-api/src/api/providers/items"
	"github.com/sirupsen/logrus"
)

type itemsService struct{}

type itemsServiceInterface interface {
	SearchItems(query string, offset int) (*items.ItemSearchResponse, apierrors.ApiError)
	GetItemWithDescription(itemId string) (*items.Item, apierrors.ApiError)
	GetItemHistory(itemId string) ([]items.ItemHistory, apierrors.ApiError)
	GetItemReviews(itemId string, catalogProductId string) (*items.ItemReviewsResponse, apierrors.ApiError)
	GetCategoryTrends(categoryId string) (*items.CategoryTrends, apierrors.ApiError)
	GetItem(itemId string) (*items.Item, apierrors.ApiError)
}

var (
	ItemsService itemsServiceInterface
)

func init() {
	ItemsService = &itemsService{}
}

func (s *itemsService) SearchItems(query string, offset int) (*items.ItemSearchResponse, apierrors.ApiError) {
	result, err := items_provider.SearchItemsByQuery(query, offset)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *itemsService) GetItemWithDescription(itemId string) (*items.Item, apierrors.ApiError) {
	var meliItem *items.Item
	var desc *items.ItemDescription
	input := make(chan items.ItemDescriptionConcurrent, 2)
	defer close(input)

	go func(itemId string, output chan items.ItemDescriptionConcurrent) {
		item, err := items_provider.GetItemById(itemId)
		output <- items.ItemDescriptionConcurrent{
			Item:        item,
			Description: nil,
			Error:       err,
		}
	}(itemId, input)

	go func(itemId string, output chan items.ItemDescriptionConcurrent) {
		description, err := items_provider.GetItemDescription(itemId)
		output <- items.ItemDescriptionConcurrent{
			Item:        nil,
			Description: description,
			Error:       err,
		}
	}(itemId, input)

	var err apierrors.ApiError = nil
	for i := 0; i < 2; i++ {
		result := <-input
		err = result.Error

		if result.Item != nil {
			category, _ := items_provider.GetCategory(result.Item.CategoryId)
			if category != nil && len(category.PathFromRoot) > 0 {
				result.Item.RootCategory = category.PathFromRoot[0]["name"]
			} else {
				result.Item.RootCategory = "Otros"
			}

			meliItem = result.Item
		} else if result.Description != nil {
			desc = result.Description
		}
	}

	if err != nil {
		return nil, err
	}

	if desc != nil && meliItem != nil {
		meliItem.Description = desc.PlainText
	}

	return meliItem, nil
}

func (s *itemsService) GetItem(itemId string) (*items.Item, apierrors.ApiError) {
	return items_provider.GetItemById(itemId)
}

func (s *itemsService) GetItemHistory(itemId string) ([]items.ItemHistory, apierrors.ApiError) {
	return items.ItemHistoryDao.GetItemHistory(itemId)
}

func (s *itemsService) GetItemReviews(itemId string, catalogProductId string) (*items.ItemReviewsResponse, apierrors.ApiError) {
	result, err := items_provider.GetItemReviews(itemId, catalogProductId)
	if err != nil {
		logrus.Error(fmt.Sprintf("error while getting reviews for item %s", itemId), err)
		return nil, err
	}

	return result, nil
}

func (s *itemsService) GetCategoryTrends(categoryId string) (*items.CategoryTrends, apierrors.ApiError) {
	return items_provider.GetCategoryTrends(categoryId)
}