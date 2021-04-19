package lists

const (
	PrivacyTypePrivate = "private"
	PrivacyTypePublic  = "public"
)

type ListDto struct {
	Id          int64  `json:"id"`
	OwnerId     int64  `json:"owner_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	DateCreated string `json:"date_created"`
}

type Lists []ListDto
