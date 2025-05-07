package websocket

import "github.com/google/uuid"

type Action string

const (
	PredefinedActionNoop   Action = "noop"
	PredefinedActionReject Action = "reject"
	PredefinedActionData   Action = "data"
	PredefinedActionNotify Action = "notify"
)

type ControlFlag string

const (
	ControlFlagQuite = ControlFlag("quite")
	ControlFlagClose = ControlFlag("close")
)

type ConnectWrapper struct {
	ConnectID string

	Authorized   bool
	UserID       int64
	UserMetadata map[string]any
}

type MessageWrapper struct {
	Action  Action
	Payload any

	SessionID string
	MessageID string

	Authorized   bool
	UserID       int64
	UserMetadata map[string]any
}

type PingWrapper struct {
	SessionID string
	MessageID string

	Authorized   bool
	UserID       int64
	UserMetadata map[string]any
}

type CloseWrapper struct {
	CloseCode int

	SessionID string

	Authorized   bool
	UserID       int64
	UserMetadata map[string]any
}

type ResponseWrapper struct {
	Action  Action `json:"action"`
	Payload any    `json:"payload"`

	SessionID string `json:"session_id"`
	MessageID string `json:"message_id"`

	ControlFlags []ControlFlag `json:"-"`
}

func BuildResponseFrom[WRAP MessageWrapper | ConnectWrapper | PingWrapper | CloseWrapper](
	request WRAP, id string, action Action, payload any, flags ...ControlFlag,
) ResponseWrapper {
	response := ResponseWrapper{
		Action:       action,
		Payload:      payload,
		ControlFlags: flags,
	}
	switch request := any(request).(type) {
	case MessageWrapper:
		response.SessionID = request.SessionID
		response.MessageID = id
	case ConnectWrapper:
		response.SessionID = request.ConnectID
		response.MessageID = id
	case PingWrapper:
		response.SessionID = request.SessionID
		response.MessageID = id
	case CloseWrapper:
		response.SessionID = request.SessionID
		response.MessageID = id
	}
	return response
}

func BuildNoopResponse[WRAP MessageWrapper | ConnectWrapper | PingWrapper | CloseWrapper](request WRAP) ResponseWrapper {
	return BuildResponseFrom(request,
		uuid.New().String(),
		PredefinedActionNoop, "none",
		ControlFlagQuite,
	)
}

func BuildCloseResponse[WRAP MessageWrapper | ConnectWrapper | PingWrapper | CloseWrapper](
	request WRAP, action Action, payload any,
) ResponseWrapper {
	return BuildResponseFrom(request,
		uuid.New().String(),
		action, payload,
		ControlFlagClose,
	)
}

func BuildRejectResponse[WRAP MessageWrapper | ConnectWrapper | PingWrapper | CloseWrapper](request WRAP) ResponseWrapper {
	return BuildResponseFrom(request,
		uuid.New().String(),
		PredefinedActionReject, "none",
		ControlFlagQuite,
		ControlFlagClose,
	)
}
