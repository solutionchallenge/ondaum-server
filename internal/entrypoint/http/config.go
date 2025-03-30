package http

import (
	"github.com/solutionchallenge/ondaum-server/pkg/database"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/jwt"
	"github.com/solutionchallenge/ondaum-server/pkg/oauth"
)

type AppConfig struct {
	Verbose        bool            `mapstructure:"verbose"`
	HttpConfig     http.Config     `mapstructure:"http"`
	DatabaseConfig database.Config `mapstructure:"database"`
	Migration      MigrationConfig `mapstructure:"migration"`
	OAuthConfig    oauth.Config    `mapstructure:"oauth"`
	JWTConfig      jwt.Config      `mapstructure:"jwt"`
}

type MigrationConfig struct {
	Enabled bool `mapstructure:"enabled"`
}
