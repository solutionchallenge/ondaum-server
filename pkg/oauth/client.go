package oauth

type Client interface {
	GetProvider() Provider
	GetAuthURL(state string, override ...string) string
	GetUserInfo(code string) (UserInfoOutput, error)
}
