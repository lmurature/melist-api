package users

type User struct {
	Id           int64       `json:"id"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Nickname     string      `json:"nickname"`
	Email        string      `json:"email"`
}
