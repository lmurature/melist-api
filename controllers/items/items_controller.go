package items_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lmurature/melist-api/domain/apierrors"
	items_service "github.com/lmurature/melist-api/services/items"
	"net/http"
	"net/url"
)

func SearchItems(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		err := apierrors.NewBadRequestApiError("search 'query' can't be empty")
		c.JSON(err.Status(), err)
		return
	}

	// TODO: Contemplate items search paging and sorting.

	result, err := items_service.ItemsService.SearchItems(url.QueryEscape(query))
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

	item, err := items_service.ItemsService.GetItem(itemId)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, item)
}
