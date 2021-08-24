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
	input := make(chan items.ItemConcurrent, len(itemIds))
	defer close(input)

	for i := range itemIds {
		go func(id string, index int, output chan items.ItemConcurrent) {
			item, err := items_service.ItemsService.GetItem(id)
			output <- items.ItemConcurrent{
				Item:      item,
				Error:     err,
				ListIndex: index,
			}
		}(itemIds[i], i, input)
	}

	for i := 0; i < len(itemIds); i++ {
		result := <-input
		if result.Error != nil {
			logrus.Error("error executing job while getting item from Mercado Libre's API")
			continue
		}

		meliItems = append(meliItems, *result.Item)
	}

	logrus.Info(fmt.Sprintf("about to save %d items data to item history table", len(meliItems)))
	// save data into items history table

	for _, item := range meliItems {
		hist := items.ItemHistory{
			ItemId:      item.Id,
			Price:       item.Price,
			Quantity:    item.AvailableQuantity,
			Status:      item.Status,
			HasDeal:     item.HasActiveDeal(),
			DateFetched: date_utils.GetNowDateFormatted(),
		}

		result, err := items.ItemHistoryDao.InsertItemHistory(hist)
		if err == nil {
			logrus.Info(fmt.Sprintf("successfully saved item history %v", result))
		}
	}
}

func persistNotifications() {
	// TODO: ...
	// get all lists

	// for each list, get all items, analyze current data and last history, generate notifications
}
