package chat

import (
	"encoding/json"
	"slices"

	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

type ChatLLMResponseType string

const (
	ChatLLMResponseTypeAction ChatLLMResponseType = "action"
	ChatLLMResponseTypeText   ChatLLMResponseType = "text"
)

var ChatLLMResponseTypeList = []ChatLLMResponseType{
	ChatLLMResponseTypeAction,
	ChatLLMResponseTypeText,
}

type ChatLLMResponseAction = common.Feature

var ChatLLMResponseActionList = []ChatLLMResponseAction(common.SupportedFeatures)

type ChatLLMResponse struct {
	Type ChatLLMResponseType `json:"type"`
	Data string              `json:"data"`
}

func ParseChatLLMResponse(response string) (ChatLLMResponse, error) {
	var result ChatLLMResponse
	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		return ChatLLMResponse{}, utils.WrapError(err, "failed to parse chat llm response")
	}
	if !slices.Contains(ChatLLMResponseTypeList, result.Type) {
		return ChatLLMResponse{}, utils.NewError("invalid type(%v)", result.Type)
	}
	if result.Type == ChatLLMResponseTypeAction && !slices.Contains(ChatLLMResponseActionList, ChatLLMResponseAction(result.Data)) {
		return ChatLLMResponse{}, utils.NewError("invalid action(%v)", result.Data)
	}
	return result, nil
}

func IsValidChatLLMResponse(response string) bool {
	_, err := ParseChatLLMResponse(response)
	return err == nil
}
