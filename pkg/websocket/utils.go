package websocket

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

func GetWebsocketUserID[T *fiberws.Conn | *fiber.Ctx](c T) (int64, error) {
	localValue := any(nil)
	switch c := any(c).(type) {
	case *fiberws.Conn:
		localValue = c.Locals("X-Websocket-User-ID")
		if localValue == nil {
			return 0, utils.NewError("X-Websocket-User-ID is not set")
		}
	case *fiber.Ctx:
		localValue = c.Locals("X-Websocket-User-ID")
		if localValue == nil {
			return 0, utils.NewError("X-Websocket-User-ID is not set")
		}
	}
	stringValue, ok := localValue.(string)
	if !ok {
		return 0, utils.NewError("X-Websocket-User-ID is not a string")
	}
	userID, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		return 0, utils.WrapError(err, "X-Websocket-User-ID is not a valid integer")
	}
	return userID, nil
}

func GetWebsocketUserMetadata[T *fiberws.Conn | *fiber.Ctx](c T) (map[string]any, error) {
	localValue := any(nil)
	switch c := any(c).(type) {
	case *fiberws.Conn:
		localValue = c.Locals("X-Websocket-User-Metadata")
		if localValue == nil {
			return nil, utils.NewError("X-Websocket-User-Metadata is not set")
		}
	case *fiber.Ctx:
		localValue = c.Locals("X-Websocket-User-Metadata")
		if localValue == nil {
			return nil, utils.NewError("X-Websocket-User-Metadata is not set")
		}
	}
	userMetadata, ok := localValue.(map[string]any)
	if !ok {
		return nil, utils.NewError("X-Websocket-User-Metadata is not a map")
	}
	return userMetadata, nil
}

func GetWebsocketSessionID[T *fiberws.Conn | *fiber.Ctx](c T) (string, error) {
	localValue := any(nil)
	switch c := any(c).(type) {
	case *fiberws.Conn:
		localValue = c.Locals("X-Websocket-Session-ID")
		if localValue == nil {
			return "", utils.NewError("X-Websocket-Session-ID is not set")
		}
	case *fiber.Ctx:
		localValue = c.Locals("X-Websocket-Session-ID")
		if localValue == nil {
			return "", utils.NewError("X-Websocket-Session-ID is not set")
		}
	}
	stringValue, ok := localValue.(string)
	if !ok {
		return "", utils.NewError("X-Websocket-Session-ID is not a string")
	}
	return stringValue, nil
}
