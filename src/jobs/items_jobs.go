package jobs

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/domain/items"
	items_service "github.com/lmurature/melist-api/src/api/services/items"
	date_utils "github.com/lmurature/melist-api/src/api/utils/date"
	"github.com/onatm/clockwerk"
	"github.com/sirupsen/logrus"
)

var (
	ItemsJobs clockwerk.Job
)

type ItemsJobsStruct struct{}

func init() {
	ItemsJobs = &ItemsJobsStruct{}
}

func (i ItemsJobsStruct) Run() {
	// analyze lists and persist notifications
	persistNotifications()

	// items saved in lists history
	persistItemHistory()
}

func persistItemHistory() {
	// get all items from items table
	itemIds, err := items.ItemDao.GetAllItems()
	if err != nil {
		logrus.Error("error executing job while getting all items from db", err)
		return
	}

	// get items data from meli api
	meliItems := make([]items.Item, 0)
	for _, id := range itemIds {
		item, itemErr := items_service.ItemsService.GetItem(id)
		if itemErr != nil {
			continue
		}
		reviews, revErr := items_service.ItemsService.GetItemReviews(id)
		if revErr != nil {
			continue
		}
		item.ReviewsQuantity = reviews.Paging.Total
		meliItems = append(meliItems, *item)
	}

	logrus.Info(fmt.Sprintf("about to save %d items data to item history table", len(meliItems)))

	// save data into items history table
	for _, item := range meliItems {
		hist := items.ItemHistory{
			ItemId:          item.Id,
			Price:           item.Price,
			Quantity:        item.AvailableQuantity,
			Status:          item.Status,
			HasDeal:         item.HasActiveDeal(),
			ReviewsQuantity: item.ReviewsQuantity,
			DateFetched:     date_utils.GetNowDateFormatted(),
		}

		_, _ = items.ItemHistoryDao.InsertItemHistory(hist)
	}
}

func persistNotifications() {
	// TODO: ...
	// get all lists

	// for each list, get all items, analyze current data and last history, generate notifications
}
