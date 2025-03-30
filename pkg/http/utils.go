package http

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

func GetRequestID(c *fiber.Ctx) string {
	requestID := c.Get("X-Request-Id")
	if requestID == "" {
		requestID = uuid.New().String()
		c.Set("X-Request-Id", requestID)
	}
	c.SetUserContext(utils.WithValue(c.UserContext(), utils.CtxKeyForRequestID, requestID))
	return requestID
}

func GetUserID(c *fiber.Ctx) (int64, error) {
	localValue := c.Locals("X-User-ID")
	if localValue == nil {
		return 0, fmt.Errorf("X-User-ID is not set")
	}
	stringValue, ok := localValue.(string)
	if !ok {
		return 0, fmt.Errorf("X-User-ID is not a string")
	}
	userID, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("X-User-ID is not a valid integer: %w", err)
	}
	return userID, nil
}

func GetUserMetadata(c *fiber.Ctx) (map[string]any, error) {
	userMetadata, ok := c.Locals("X-User-Metadata").(map[string]any)
	if !ok {
		return nil, fmt.Errorf("X-User-Metadata is not set")
	}
	return userMetadata, nil
}
