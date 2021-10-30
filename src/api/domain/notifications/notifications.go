package notifications

import (
	"fmt"
	date_utils "github.com/lmurature/melist-api/src/api/utils/date"
)

const (
	listItemUrl = "/lists/%d/%s"
	listUrl     = "/lists/%d"
	listItemReviewsUrl = "/lists/%d/%s/reviews"
)

type Notification struct {
	Id        int64  `json:"id"`
	ListId    int64  `json:"list_id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Permalink string `json:"permalink"m `
	Seen      bool   `json:"seen"`
}

func NewPriceChangeNotification(listId int64, itemId string, oldPrice float32, newPrice float32, title string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("¡El producto %s tuvo un cambio en su precio! Antes valía %.2f, ahora %.2f.", title, oldPrice, newPrice),
		Timestamp: date_utils.GetNowDateFormatted(),
		Permalink: fmt.Sprintf(listItemUrl, listId, itemId),
	}
}

func NewDealActivatedNotification(listId int64, itemId string, title string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("¡El producto %s entró en una oferta!", title),
		Timestamp: date_utils.GetNowDateFormatted(),
		Permalink: fmt.Sprintf(listItemUrl, listId, itemId),
	}
}

func NewDealEndedNotification(listId int64, itemId string, title string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("El producto %s ya no está en oferta.", title),
		Timestamp: date_utils.GetNowDateFormatted(),
		Permalink: fmt.Sprintf(listItemUrl, listId, itemId),
	}
}

func NewNearEmptyStockNotification(listId int64, itemId string, currentStock int, title string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("El producto %s se está por quedar sin stock, sólo restan %d unidades disponibles.", title, currentStock),
		Timestamp: date_utils.GetNowDateFormatted(),
		Permalink: fmt.Sprintf(listItemUrl, listId, itemId),
	}
}

func NewItemChangedStatusNotification(listId int64, itemId string, title string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("El producto %s ya no se puede comprar. La publicación fue pausada o finalizada.", title),
		Timestamp: date_utils.GetNowDateFormatted(),
		Permalink: fmt.Sprintf(listItemUrl, listId, itemId),
	}
}

func NewEmptyStockNotification(listId int64, itemId string, title string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("El producto %s se quedó sin stock.", title),
		Timestamp: date_utils.GetNowDateFormatted(),
		Permalink: fmt.Sprintf(listItemUrl, listId, itemId),
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

func NewUncheckedItemNotification(listId int64, itemId string, checkerUser string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("El producto %s fue marcado como no comprado por %s", itemId, checkerUser),
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

func NewReviewItemNotification(listId int64, itemId string, title string) *Notification {
	return &Notification{
		ListId:    listId,
		Message:   fmt.Sprintf("¡El producto %s tiene revisiones nuevas por parte de otros usuarios de Mercado Libre!", title),
		Timestamp: date_utils.GetNowDateFormatted(),
		Permalink: fmt.Sprintf(listItemReviewsUrl, listId, itemId),
	}
}
