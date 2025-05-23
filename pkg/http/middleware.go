package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/pkg/jwt"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

type MiddlewareFunc = func(c *fiber.Ctx) error

func NewJWTAuthMiddleware(generator jwt.Generator) MiddlewareFunc {
	return func(c *fiber.Ctx) error {
		rid := GetRequestID(c)

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			utils.Log(utils.DebugLevel).Ctx(c.UserContext()).RID(rid).BT().Send("No authorization header")
			return c.Next()
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			utils.Log(utils.WarnLevel).Ctx(c.UserContext()).RID(rid).BT().Send("Missing Bearer prefix in Authorization header")
			return c.Next()
		}

		tokenString := authHeader[7:]
		tokenType, err := generator.GetTokenType(tokenString)
		if err != nil {
			utils.Log(utils.WarnLevel).Ctx(c.UserContext()).Err(err).RID(rid).BT().Send("Failed to get token type")
			return c.Next()
		}

		if tokenType != jwt.AccessTokenType {
			utils.Log(utils.WarnLevel).Ctx(c.UserContext()).RID(rid).BT().Send("Invalid token type")
			return c.Next()
		}

		claims, err := generator.UnpackToken(tokenString)
		if err != nil {
			utils.Log(utils.InfoLevel).Ctx(c.UserContext()).Err(err).RID(rid).BT().Send("Failed to unpack access token")
			return c.Next()
		}

		utils.Log(utils.InfoLevel).Ctx(c.UserContext()).RID(rid).BT().Send("Unpacked User: %v (%+v)", claims.Value, claims.Metadata)
		c.Locals("X-User-Metadata", claims.Metadata)
		c.Locals("X-User-ID", claims.Value)

		return c.Next()
	}
}
