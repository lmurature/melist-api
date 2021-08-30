package jobs

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/domain/items"
	"github.com/lmurature/melist-api/src/api/domain/notifications"
	items_provider "github.com/lmurature/melist-api/src/api/providers/items"
	items_service "github.com/lmurature/melist-api/src/api/services/items"
	lists_service "github.com/lmurature/melist-api/src/api/services/lists"
	date_utils "github.com/lmurature/melist-api/src/api/utils/date"
	"github.com/onatm/clockwerk"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	ItemsJobs clockwerk.Job
)

type ItemsJobsStruct struct{}

func init() {
	ItemsJobs = &ItemsJobsStruct{}
}

func (i ItemsJobsStruct) Run() {
	// analyze lists, persist history and notifications
	persistNotifications()
}

func persistNotifications() {
	// get all lists
	lists, err := lists_service.ListsService.GetAllLists()
	if err != nil {
		logrus.Error("error while getting lists", err)
		return
	}

	// for each list, get all items, analyze current data and last history, generate notifications and save item history.
	for _, list := range lists {
		listItems, err := lists_service.ListsService.GetItemsFromList(list.Id, list.OwnerId, true)
		if err != nil {
			logrus.Error("error while getting list items")
			continue
		}

		logrus.Info(fmt.Sprintf("about to analyze and save history for %d items from list %d", len(listItems), list.Id))
		for _, item := range listItems {
			lastHistory, err := items.ItemHistoryDao.GetLastItemHistory(item.ItemId)
			if err != nil && err.Status() != http.StatusNotFound {
				logrus.Error("error getting item last history", err)
				continue
			}

			reviews, revErr := items_service.ItemsService.GetItemReviews(item.MeliItem.Id, item.MeliItem.CatalogProductId)
			if revErr != nil {
				continue
			}

			realQ, err := items_provider.GetRealQuantity(item.MeliItem.Permalink)
			if err == nil {
				item.MeliItem.AvailableQuantity = int(*realQ)
			}


			item.MeliItem.ReviewsQuantity = reviews.Paging.Total

			// analyze
			if lastHistory != nil {
				if item.MeliItem.HasActiveDeal() && !lastHistory.HasDeal {
					_, _ = notifications.NotificationsDao.SaveNotification(*notifications.NewDealActivatedNotification(list.Id, item.ItemId))
				} else if !item.MeliItem.HasActiveDeal() && lastHistory.HasDeal {
					_, _ = notifications.NotificationsDao.SaveNotification(*notifications.NewDealEndedNotification(list.Id, item.ItemId))
				}

				if item.MeliItem.Price != lastHistory.Price {
					_, _ = notifications.NotificationsDao.SaveNotification(*notifications.NewPriceChangeNotification(list.Id, item.ItemId, lastHistory.Price, item.MeliItem.Price))
				}

				if item.MeliItem.AvailableQuantity == 0 &&
					item.MeliItem.Status == "paused" &&
					lastHistory.Quantity > 0 {
					_, _ = notifications.NotificationsDao.SaveNotification(*notifications.NewEmptyStockNotification(list.Id, item.ItemId))
				}

				if item.MeliItem.AvailableQuantity <= 3 && lastHistory.Quantity > 3 {
					_, _ = notifications.NotificationsDao.SaveNotification(*notifications.NewNearEmptyStockNotification(list.Id, item.ItemId, item.MeliItem.AvailableQuantity))
				}

				if item.MeliItem.ReviewsQuantity > lastHistory.ReviewsQuantity {
					_, _ = notifications.NotificationsDao.SaveNotification(*notifications.NewReviewItemNotification(list.Id, item.ItemId))
				}
			}

			hist := items.ItemHistory{
				ItemId:          item.MeliItem.Id,
				Price:           item.MeliItem.Price,
				Quantity:        item.MeliItem.AvailableQuantity,
				Status:          item.MeliItem.Status,
				HasDeal:         item.MeliItem.HasActiveDeal(),
				ReviewsQuantity: item.MeliItem.ReviewsQuantity,
				DateFetched:     date_utils.GetNowDateFormatted(),
			}

			_, _ = items.ItemHistoryDao.InsertItemHistory(hist)
		}
	}
}
