package websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/solutionchallenge/ondaum-server/pkg/jwt"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

type CoreMessage struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

func EnableWebsocketCore(app fiber.Router, path string, generator *jwt.Generator) error {
	app.Use(path, func(c *fiber.Ctx) error {
		if c.Get("Upgrade") != "websocket" {
			return c.SendStatus(fiber.StatusUpgradeRequired)
		}
		sid := c.Query("session_id")
		if sid == "" {
			sid = uuid.New().String()
		}
		c.Locals("X-Websocket-Session-ID", sid)
		tryWebsocketAuthorization(c, generator, sid)
		return c.Next()
	})
	return nil
}

func tryWebsocketAuthorization(c *fiber.Ctx, generator *jwt.Generator, sessionID string) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		utils.Log(utils.DebugLevel).Ctx(c.UserContext()).RID(sessionID).BT().Send("No authorization header")
		return
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		utils.Log(utils.WarnLevel).Ctx(c.UserContext()).RID(sessionID).BT().Send("Missing Bearer prefix in Authorization header")
		return
	}

	tokenString := authHeader[7:]
	tokenType, err := generator.GetTokenType(tokenString)
	if err != nil {
		utils.Log(utils.WarnLevel).Ctx(c.UserContext()).Err(err).RID(sessionID).BT().Send("Failed to get token type")
		return
	}

	if tokenType != jwt.AccessTokenType {
		utils.Log(utils.WarnLevel).Ctx(c.UserContext()).RID(sessionID).BT().Send("Required valid access token")
		return
	}

	claims, err := generator.UnpackToken(tokenString)
	if err != nil {
		utils.Log(utils.WarnLevel).Ctx(c.UserContext()).Err(err).RID(sessionID).BT().Send("Failed to unpack token")
		return
	}

	utils.Log(utils.InfoLevel).Ctx(c.UserContext()).RID(sessionID).BT().Send("Unpacked User: %v (%+v)", claims.Value, claims.Metadata)
	c.Locals("X-Websocket-User-Metadata", claims.Metadata)
	c.Locals("X-Websocket-User-ID", claims.Value)
}
