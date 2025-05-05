package oauth

type Client interface {
	GetProvider() Provider
	GetAuthURL(state string, override ...string) (string, error)
	GetUserInfo(code string, override ...string) (UserInfoOutput, error)
}
