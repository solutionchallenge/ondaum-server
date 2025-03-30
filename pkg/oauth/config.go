package oauth

type Config struct {
	Google struct {
		Enabled      bool   `mapstructure:"enabled"`
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		RedirectURL  string `mapstructure:"redirect_url"`
	} `mapstructure:"google"`
}
