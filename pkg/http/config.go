package http

import "time"

type Config struct {
	Host  string `mapstructure:"host" default:"0.0.0.0"`
	Port  int    `mapstructure:"port" default:"8080"`
	Limit struct {
		Payload     int `mapstructure:"payload" default:"1024"`
		Concurrency int `mapstructure:"concurrency" default:"100"`
	} `mapstructure:"limit"`
	Buffer struct {
		Read  int `mapstructure:"read" default:"4096"`
		Write int `mapstructure:"write" default:"4096"`
	} `mapstructure:"buffer"`
	Timeout struct {
		Read     time.Duration `mapstructure:"read" default:"10s"`
		Write    time.Duration `mapstructure:"write" default:"10s"`
		Idle     time.Duration `mapstructure:"idle" default:"10s"`
		Shutdown time.Duration `mapstructure:"shutdown" default:"10s"`
	} `mapstructure:"timeout"`
}
