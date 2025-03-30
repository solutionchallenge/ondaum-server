package jwt

type Config struct {
	SecretKey     string `mapstructure:"secret_key"`
	AccessExpire  int    `mapstructure:"access_expire"`
	RefreshExpire int    `mapstructure:"refresh_expire"`
}
