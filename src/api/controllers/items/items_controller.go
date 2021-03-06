package items_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	items_service "github.com/lmurature/melist-api/src/api/services/items"
	"net/http"
	"net/url"
	"strconv"
)

func SearchItems(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		err := apierrors.NewBadRequestApiError("search 'query' can't be empty")
		c.JSON(err.Status(), err)
		return
	}

	offsetParam := c.DefaultQuery("offset", "0")
	offset, parseErr := strconv.Atoi(offsetParam)
	if parseErr != nil {
		br := apierrors.NewBadRequestApiError("offset must be a number")
		c.JSON(br.Status(), br)
		return
	}

	result, err := items_service.ItemsService.SearchItems(url.QueryEscape(query), offset)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetItem(c *gin.Context) {
	itemId := c.Param("item_id")
	if itemId == "" {
		err := apierrors.NewBadRequestApiError("'item_id' can't be empty")
		c.JSON(err.Status(), err)
		return
	}

	item, err := items_service.ItemsService.GetItemWithDescription(itemId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func GetItemHistory(c *gin.Context) {
	itemId := c.Param("item_id")
	if itemId == "" {
		err := apierrors.NewBadRequestApiError("'item_id' can't be empty")
		c.JSON(err.Status(), err)
		return
	}

	result, err := items_service.ItemsService.GetItemHistory(itemId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetItemReviews(c *gin.Context) {
	itemId := c.Param("item_id")
	if itemId == "" {
		err := apierrors.NewBadRequestApiError("'item_id' can't be empty")
		c.JSON(err.Status(), err)
		return
	}

	result, err := items_service.ItemsService.GetItemReviews(itemId, c.Query("catalog_product_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetCategoryTrends(c *gin.Context) {
	categoryId := c.Param("category_id")
	if categoryId == "" {
		err := apierrors.NewBadRequestApiError("'category_id' can't be empty")
		c.JSON(err.Status(), err)
		return
	}

	result, err := items_service.ItemsService.GetCategoryTrends(categoryId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result)
}
