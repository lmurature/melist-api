package items_controller

import (
	"github.com/gin-gonic/gin"
	apierrors2 "github.com/lmurature/melist-api/src/api/domain/apierrors"
	items_service2 "github.com/lmurature/melist-api/src/api/services/items"
	"net/http"
	"net/url"
)

func SearchItems(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		err := apierrors2.NewBadRequestApiError("search 'query' can't be empty")
		c.JSON(err.Status(), err)
		return
	}

	// TODO: Contemplate items search paging and sorting.

	result, err := items_service2.ItemsService.SearchItems(url.QueryEscape(query))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetItem(c *gin.Context) {
	itemId := c.Param("item_id")
	if itemId == "" {
		err := apierrors2.NewBadRequestApiError("'item_id' can't be empty")
		c.JSON(err.Status(), err)
		return
	}

	item, err := items_service2.ItemsService.GetItem(itemId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, item)
}
