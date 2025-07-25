package authentication

// Token subject constants for distinguishing token types
const (
	AccessTokenSubject  = "access_token"
	RefreshTokenSubject = "refresh_token"
)

// Session represents the output of a successful authentication (login), containing issued tokens
type Session struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
