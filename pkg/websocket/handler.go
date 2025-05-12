package websocket

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

type Handler interface {
	Identify() string
	HandleConnect(c *fiberws.Conn, payload ConnectWrapper) (ResponseWrapper, string, error)
	HandleMessage(c *fiberws.Conn, payload MessageWrapper) (ResponseWrapper, bool, error)
	HandlePing(c *fiberws.Conn, payload PingWrapper) (ResponseWrapper, bool, error)
	HandleClose(c *fiberws.Conn, payload CloseWrapper)
}

func Install(router fiber.Router, path string, handler Handler) fiber.Router {
	return router.Get(path, fiberws.New(func(c *fiberws.Conn) {
		responseWrapper, sessionID, err := handleConnect(c, handler)
		if err != nil {
			utils.Log(utils.ErrorLevel).Err(err).CID(sessionID).BT().Send("Failed to connect")
			closeConnection(c, sessionID, "", "failed to connect", fiberws.CloseProtocolError)
			return
		}
		closed, err := processControlFlags(c, responseWrapper)
		if err != nil {
			utils.Log(utils.ErrorLevel).Err(err).CID(sessionID).BT().Send("Failed to process control flags")
			closeConnection(c, sessionID, "", "failed to process control flags", fiberws.CloseProtocolError)
			return
		} else if closed {
			utils.Log(utils.InfoLevel).CID(sessionID).BT().Send("Closing connection")
			closeConnection(c, sessionID, "", "server requested")
			return
		}
		for {
			messageID := uuid.New().String()
			messageType, rawMessage, err := c.ReadMessage()
			if err != nil {
				switch {
				case fiberws.IsCloseError(err, fiberws.CloseNormalClosure, fiberws.CloseGoingAway, fiberws.CloseAbnormalClosure):
					utils.Log(utils.InfoLevel).Err(err).CID(sessionID).RID(messageID).BT().Send("Connection closed by client")
					handleClose(c, sessionID, fiberws.CloseNormalClosure, handler)
					return
				default:
					utils.Log(utils.ErrorLevel).Err(err).CID(sessionID).RID(messageID).BT().Send("Failed to read message")
					handleClose(c, sessionID, fiberws.CloseInternalServerErr, handler)
					return
				}
			}
			switch messageType {
			case fiberws.TextMessage, fiberws.BinaryMessage:
				responseWrapper, isCritical, err := handleMessage(c, sessionID, messageID, rawMessage, handler)
				if err != nil {
					if isCritical {
						utils.Log(utils.ErrorLevel).Err(err).CID(sessionID).RID(messageID).BT().Send("Occurred critical error")
						closeConnection(c, sessionID, messageID, "server occurred critical error", fiberws.CloseInternalServerErr)
						handleClose(c, sessionID, fiberws.CloseInternalServerErr, handler)
						return
					}
					utils.Log(utils.ErrorLevel).Err(err).CID(sessionID).RID(messageID).BT().Send("Failed to handle payload message")
					continue
				}
				closed, err := processControlFlags(c, responseWrapper)
				if err != nil {
					utils.Log(utils.ErrorLevel).Err(err).CID(sessionID).RID(messageID).BT().Send("Failed to process control flags")
					continue
				} else if closed {
					utils.Log(utils.InfoLevel).CID(sessionID).RID(messageID).BT().Send("Closing connection")
					closeConnection(c, sessionID, messageID, "server requested")
					return
				}
			case fiberws.PingMessage:
				responseWrapper, isCritical, err := handlePing(c, sessionID, messageID, handler)
				if err != nil {
					if isCritical {
						utils.Log(utils.ErrorLevel).Err(err).CID(sessionID).RID(messageID).BT().Send("Occurred critical error in ping")
						closeConnection(c, sessionID, messageID, "server occurred critical error", fiberws.CloseInternalServerErr)
						handleClose(c, sessionID, fiberws.CloseInternalServerErr, handler)
						return
					}
					utils.Log(utils.ErrorLevel).Err(err).CID(sessionID).RID(messageID).BT().Send("Failed to handle ping")
					continue
				}
				closed, err := processControlFlags(c, responseWrapper)
				if err != nil {
					utils.Log(utils.ErrorLevel).Err(err).CID(sessionID).RID(messageID).BT().Send("Failed to process control flags")
					continue
				} else if closed {
					utils.Log(utils.InfoLevel).CID(sessionID).RID(messageID).BT().Send("Closing connection")
					closeConnection(c, sessionID, messageID, "server requested")
					return
				}
			case fiberws.PongMessage:
				continue
			case fiberws.CloseMessage:
				utils.Log(utils.InfoLevel).CID(sessionID).RID(messageID).BT().Send("Connection closed")
				handleClose(c, sessionID, fiberws.CloseNormalClosure, handler)
				// already closed by client so no need to handle
				return
			default:
				utils.Log(utils.WarnLevel).CID(sessionID).RID(messageID).BT().Send("Unknown message type")
			}
		}
	})).Name(handler.Identify())
}

func handleConnect(
	c *fiberws.Conn, handler Handler,
) (ResponseWrapper, string, error) {
	sessionID, err := GetWebsocketSessionID(c)
	if err != nil {
		return ResponseWrapper{}, "", utils.WrapError(err, "failed to get session ID")
	}
	connectWrapper := ConnectWrapper{
		ConnectID: sessionID,
	}
	userID, err := GetWebsocketUserID(c)
	if err == nil {
		connectWrapper.Authorized = true
		connectWrapper.UserID = userID
		userMetadata, err := GetWebsocketUserMetadata(c)
		if err == nil {
			connectWrapper.UserMetadata = userMetadata
		}
	}
	return handler.HandleConnect(c, connectWrapper)
}

func handleMessage(
	c *fiberws.Conn, sessionID string, messageID string, rawMessage []byte, handler Handler,
) (ResponseWrapper, bool, error) {
	var requestWrapper MessageWrapper
	err := json.Unmarshal(rawMessage, &requestWrapper)
	if err != nil {
		return ResponseWrapper{}, false, utils.WrapError(err, "failed to unmarshal message")
	}
	requestWrapper.MessageID = messageID
	requestWrapper.SessionID = sessionID
	userID, err := GetWebsocketUserID(c)
	if err == nil {
		requestWrapper.Authorized = true
		requestWrapper.UserID = userID
		userMetadata, err := GetWebsocketUserMetadata(c)
		if err == nil {
			requestWrapper.UserMetadata = userMetadata
		}
	}
	return handler.HandleMessage(c, requestWrapper)
}

func handlePing(
	c *fiberws.Conn, sessionID string, messageID string, handler Handler,
) (ResponseWrapper, bool, error) {
	pingWrapper := PingWrapper{
		SessionID: sessionID,
		MessageID: messageID,
	}
	userID, err := GetWebsocketUserID(c)
	if err == nil {
		pingWrapper.Authorized = true
		pingWrapper.UserID = userID
		userMetadata, err := GetWebsocketUserMetadata(c)
		if err == nil {
			pingWrapper.UserMetadata = userMetadata
		}
	}
	return handler.HandlePing(c, pingWrapper)
}

func handleClose(
	c *fiberws.Conn, sessionID string, closeCode int, handler Handler,
) {
	closeWrapper := CloseWrapper{
		CloseCode: closeCode,
		SessionID: sessionID,
	}
	userID, err := GetWebsocketUserID(c)
	if err == nil {
		closeWrapper.Authorized = true
		closeWrapper.UserID = userID
		userMetadata, err := GetWebsocketUserMetadata(c)
		if err == nil {
			closeWrapper.UserMetadata = userMetadata
		}
	}
	handler.HandleClose(c, closeWrapper)
}

func processControlFlags(
	c *fiberws.Conn,
	responseWrapper ResponseWrapper,
) (bool, error) {
	if !slices.Contains(responseWrapper.ControlFlags, ControlFlagQuite) {
		serialized, err := json.Marshal(responseWrapper)
		if err != nil {
			return false, utils.WrapError(err, "failed to serialize message")
		}
		err = c.WriteMessage(fiberws.TextMessage, serialized)
		if err != nil {
			return false, utils.WrapError(err, "failed to write message")
		}
	}
	if slices.Contains(responseWrapper.ControlFlags, ControlFlagClose) {
		return true, nil
	}
	return false, nil
}

func closeConnection(c *fiberws.Conn, sessionID string, messageID string, reason string, cause ...int) {
	code := fiberws.CloseNormalClosure
	if len(cause) > 0 {
		code = cause[0]
	}
	payload := fmt.Sprintf("connection closed by server: %s", reason)
	message := fiberws.FormatCloseMessage(code, payload)
	if err := c.WriteControl(fiberws.CloseMessage, message, time.Now().UTC().Add(time.Second)); err != nil {
		utils.Log(utils.ErrorLevel).Err(err).CID(sessionID).RID(messageID).BT().Send("Failed to send close message")
	}
	c.Close()
}
