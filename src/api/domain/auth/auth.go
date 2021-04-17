package auth

const (
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypeRefreshToken      = "refresh_token"
)

type MeliAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	UserId       int64 `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

type MeliAuthRequest struct {
	GrantType    string `json:"grant_type"`
	ClientId     int64  `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	RedirectUri  string `json:"redirect_uri,omitempty"`
}

type ClientAuthRequest struct {
	AuthorizationCode string `json:"authorization_code,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

