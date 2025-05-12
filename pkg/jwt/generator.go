package jwt

import (
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

var (
	ErrInvalidKey     = jwt.ErrInvalidKey
	ErrInvalidKeyType = jwt.ErrInvalidKeyType
	ErrTokenExpired   = jwt.ErrTokenExpired
)

type Generator struct {
	Config Config
	Clock  clock.Clock
}

func NewGenerator(config Config, clk clock.Clock) *Generator {
	return &Generator{
		Config: config,
		Clock:  clk,
	}
}

func (g *Generator) GenerateTokenPair(value string, metadata ...map[string]any) (*TokenPair, error) {
	injector := map[string]any{}
	if len(metadata) > 0 && metadata[0] != nil {
		injector = metadata[0]
	}

	accessToken, err := g.GenerateToken(
		AccessTokenType,
		value,
		injector,
		time.Duration(g.Config.AccessExpire)*time.Second,
	)
	if err != nil {
		return nil, utils.WrapError(err, "failed to generate access token")
	}

	refreshToken, err := g.GenerateToken(
		RefreshTokenType,
		value,
		injector,
		time.Duration(g.Config.RefreshExpire)*time.Second,
	)
	if err != nil {
		return nil, utils.WrapError(err, "failed to generate refresh token")
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (g *Generator) GenerateToken(typ Type, value string, metadata map[string]any, duration time.Duration) (string, error) {
	now := g.Clock.Now().UTC().UTC()
	claims := Claims{
		Type:     typ,
		Value:    value,
		Metadata: metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(g.Config.SecretKey))
}

func (g *Generator) UnpackToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(g.Config.SecretKey), nil
	})

	if err != nil {
		return nil, utils.WrapError(err, "failed to unpack token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	if claims.ExpiresAt.Before(g.Clock.Now().UTC().UTC()) {
		return nil, jwt.ErrTokenExpired
	}

	return claims, nil
}

func (g *Generator) RefreshTokenPair(refreshTokenString string) (*TokenPair, error) {
	claims, err := g.UnpackToken(refreshTokenString)
	if err != nil {
		return nil, utils.WrapError(err, "failed to unpack refresh token")
	}

	if claims.Type != RefreshTokenType {
		return nil, jwt.ErrInvalidKeyType
	}

	accessToken, err := g.GenerateToken(
		AccessTokenType,
		claims.Value,
		claims.Metadata,
		time.Duration(g.Config.AccessExpire)*time.Second,
	)
	if err != nil {
		return nil, utils.WrapError(err, "failed to generate access token")
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
	}, nil
}

func (g *Generator) GetTokenType(tokenString string) (Type, error) {
	claims, err := g.UnpackToken(tokenString)
	if err != nil {
		return InvalidType, utils.WrapError(err, "failed to get token type")
	}
	return claims.Type, nil
}
