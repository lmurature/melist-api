package items

type ItemHistory struct {
	Id              int64   `json:"id"`
	ItemId          string  `json:"item_id"`
	Price           float32 `json:"price"`
	Quantity        int     `json:"quantity"`
	Status          string  `json:"status"`
	HasDeal         bool    `json:"has_deal"`
	DateFetched     string  `json:"date_fetched"`
	ReviewsQuantity int64   `json:"reviews_quantity"`
}
