package items

type ItemReviewsResponse struct {
	Paging    Paging   `json:"paging"`
	Reviews   []Review `json:"reviews"`
	RatingAvg float32  `json:"rating_average"`
}

type Review struct {
	ReviewId     string `json:"id"`
	DateCreated  string `json:"date_created"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Rate         int    `json:"rate"`
	Valorization int    `json:"valorization"`
	Likes        int    `json:"likes"`
	Dislikes     int    `json:"dislikes"`
	BuyingDate   string `json:"buying_date"`
}
