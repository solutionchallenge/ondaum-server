package websocket

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
)

func GetWebsocketUserID[T *fiberws.Conn | *fiber.Ctx](c T) (int64, error) {
	localValue := any(nil)
	switch c := any(c).(type) {
	case *fiberws.Conn:
		localValue = c.Locals("X-Websocket-User-ID")
		if localValue == nil {
			return 0, fmt.Errorf("X-Websocket-User-ID is not set")
		}
	case *fiber.Ctx:
		localValue = c.Locals("X-Websocket-User-ID")
		if localValue == nil {
			return 0, fmt.Errorf("X-Websocket-User-ID is not set")
		}
	}
	stringValue, ok := localValue.(string)
	if !ok {
		return 0, fmt.Errorf("X-Websocket-User-ID is not a string")
	}
	userID, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("X-Websocket-User-ID is not a valid integer: %w", err)
	}
	return userID, nil
}

func GetWebsocketUserMetadata[T *fiberws.Conn | *fiber.Ctx](c T) (map[string]any, error) {
	localValue := any(nil)
	switch c := any(c).(type) {
	case *fiberws.Conn:
		localValue = c.Locals("X-Websocket-User-Metadata")
		if localValue == nil {
			return nil, fmt.Errorf("X-Websocket-User-Metadata is not set")
		}
	case *fiber.Ctx:
		localValue = c.Locals("X-Websocket-User-Metadata")
		if localValue == nil {
			return nil, fmt.Errorf("X-Websocket-User-Metadata is not set")
		}
	}
	userMetadata, ok := localValue.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("X-Websocket-User-Metadata is not a map")
	}
	return userMetadata, nil
}

func GetWebsocketSessionID[T *fiberws.Conn | *fiber.Ctx](c T) (string, error) {
	localValue := any(nil)
	switch c := any(c).(type) {
	case *fiberws.Conn:
		localValue = c.Locals("X-Websocket-Session-ID")
		if localValue == nil {
			return "", fmt.Errorf("X-Websocket-Session-ID is not set")
		}
	case *fiber.Ctx:
		localValue = c.Locals("X-Websocket-Session-ID")
		if localValue == nil {
			return "", fmt.Errorf("X-Websocket-Session-ID is not set")
		}
	}
	stringValue, ok := localValue.(string)
	if !ok {
		return "", fmt.Errorf("X-Websocket-Session-ID is not a string")
	}
	return stringValue, nil
}
