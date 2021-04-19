package share

const (
	ShareTypeRead = "read"
	ShareTypeWrite = "write"
	ShareTypeCheck = "check"
)

type ShareConfig struct {
	ListId    int64  `json:"list_id"`
	UserId    int64  `json:"user_id"`
	ShareType string `json:"share_type"`
}
