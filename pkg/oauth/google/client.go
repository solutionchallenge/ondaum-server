package google

import (
	"context"
	"encoding/json"
	"slices"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/solutionchallenge/ondaum-server/pkg/oauth"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

const (
	Provider = oauth.Provider("google")
)

type Client struct {
	BaseConfig oauth.Config
	coreConfig oauth2.Config
}

func NewClient(config oauth.Config) oauth.Client {
	if !config.Google.Enabled {
		return nil
	}

	oauthConfig := oauth2.Config{
		ClientID:     config.Google.ClientID,
		ClientSecret: config.Google.ClientSecret,
		RedirectURL:  config.Google.DefaultRedirection,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &Client{
		BaseConfig: config,
		coreConfig: oauthConfig,
	}
}

func (client *Client) GetProvider() oauth.Provider {
	return Provider
}

func (client *Client) GetAuthURL(state string, override ...string) (string, error) {
	if len(override) > 0 && override[0] != "" {
		if !slices.Contains(client.BaseConfig.Google.AllowedRedirections, override[0]) {
			return "", utils.NewError("invalid override redirect URL: %s", override[0])
		}
		copied := client.coreConfig
		copied.RedirectURL = override[0]
		return copied.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
	}
	return client.coreConfig.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (client *Client) GetUserInfo(code string, override ...string) (oauth.UserInfoOutput, error) {
	token := (*oauth2.Token)(nil)
	err := error(nil)
	var current oauth2.Config
	if len(override) > 0 && override[0] != "" {
		if !slices.Contains(client.BaseConfig.Google.AllowedRedirections, override[0]) {
			return oauth.UserInfoOutput{}, utils.NewError("invalid override redirect URL: %s", override[0])
		}
		copied := client.coreConfig
		copied.RedirectURL = override[0]
		token, err = copied.Exchange(context.Background(), code)
		current = copied
	} else {
		token, err = client.coreConfig.Exchange(context.Background(), code)
		current = client.coreConfig
	}
	if err != nil || token == nil {
		return oauth.UserInfoOutput{}, utils.NewError("failed to exchange token: %v", err)
	}

	core := current.Client(context.Background(), token)
	resp, err := core.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return oauth.UserInfoOutput{}, utils.NewError("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	var userInfo oauth.UserInfoOutput
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return oauth.UserInfoOutput{}, utils.NewError("failed to decode user info: %v", err)
	}

	return userInfo, nil
}
