package notifications

import (
	"fmt"
	date_utils "github.com/lmurature/melist-api/src/api/utils/date"
)

// TODO: change item id into titles?

const (
	listItemUrl = "/lists/%d/%s"
	listUrl = "/lists/%d"
)

type Notification struct {
	Id        int64  `json:"id"`
	ListId    int64  `json:"list_id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Permalink string `json:"permalink"`
	Seen      bool   `json:"seen"`
}

func NewPriceChangeNotification(listId int64, itemId string, oldPrice float32, newPrice float32) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("¡El producto %s tuvo un cambio en su precio! Antes valía %f, ahora %f.", itemId, oldPrice, newPrice),
		Timestamp: date_utils.GetNowDateFormatted(),
	}
}

func NewDealActivatedNotification(listId int64, itemId string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("¡El producto %s entró en una oferta!", itemId),
		Timestamp: date_utils.GetNowDateFormatted(),
	}
}

func NewDealEndedNotification(listId int64, itemId string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("El producto %s ya no está en oferta.", itemId),
		Timestamp: date_utils.GetNowDateFormatted(),
	}
}

func NewNearEmptyStockNotification(listId int64, itemId string, currentStock int) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("El producto %s se está por quedar sin stock, sólo restan %d unidades disponibles.", itemId, currentStock),
		Timestamp: date_utils.GetNowDateFormatted(),
	}
}

func NewEmptyStockNotification(listId int64, itemId string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("El producto %s se quedó sin stock.", itemId),
		Timestamp: date_utils.GetNowDateFormatted(),
	}
}

func NewCheckedItemNotification(listId int64, itemId string, checkerUser string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("¡El producto %s fue comprado por %s!", itemId, checkerUser),
		Timestamp: date_utils.GetNowDateFormatted(),
		Permalink: fmt.Sprintf(listItemUrl, listId, itemId),
	}
}

func NewAddedItemToListNotification(listId int64, itemId string, adderUser string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("¡%s añadió un nuevo producto a la lista!", adderUser),
		Timestamp: date_utils.GetNowDateFormatted(),
		Permalink: fmt.Sprintf(listItemUrl, listId, itemId),
	}
}

func NewUserAddedListToFavorites(listId int64, favoriteUser string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("¡%s añadió esta lista a sus favoritos!", favoriteUser),
		Timestamp: date_utils.GetNowDateFormatted(),
		Permalink: fmt.Sprintf(listUrl, listId),
	}
}
