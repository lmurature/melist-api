package users

type UserDto struct {
	Id           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Nickname     string `json:"nickname"`
	Email        string `json:"email"`
	DateCreated  string `json:"date_created"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
