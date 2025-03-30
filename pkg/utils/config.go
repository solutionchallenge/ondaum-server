package utils

import (
	"strings"

	"github.com/spf13/viper"
)

func LoadConfigTo[T any](ref *T, filename string, path ...string) {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	for _, p := range path {
		viper.AddConfigPath(p)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		panic(WrapError(err, "Failed to read config file"))
	}
	err = viper.Unmarshal(ref)
	if err != nil {
		panic(WrapError(err, "Failed to unmarshal config"))
	}
}

func LoadConfig(filename string, path ...string) map[string]any {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	for _, p := range path {
		viper.AddConfigPath(p)
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(WrapError(err, "Failed to read config file"))
	}
	return viper.AllSettings()
}
