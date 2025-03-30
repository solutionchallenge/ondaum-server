package jwt

import "github.com/golang-jwt/jwt/v5"

type Type string

const (
	InvalidType      Type = "invalid"
	AccessTokenType  Type = "access_token"
	RefreshTokenType Type = "refresh_token"
)

type Claims struct {
	Type     Type           `json:"type"`
	Value    string         `json:"value"`
	Metadata map[string]any `json:"metadata"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
