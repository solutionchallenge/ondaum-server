package oauth

type Config struct {
	Google struct {
		Enabled             bool     `mapstructure:"enabled"`
		ClientID            string   `mapstructure:"client_id"`
		ClientSecret        string   `mapstructure:"client_secret"`
		DefaultRedirection  string   `mapstructure:"default_redirection"`
		AllowedRedirections []string `mapstructure:"allowed_redirections"`
	} `mapstructure:"google"`
}
