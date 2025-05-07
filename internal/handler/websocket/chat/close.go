package chat

import (
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	wspkg "github.com/solutionchallenge/ondaum-server/pkg/websocket"
	"github.com/uptrace/bun"
)

func HandleClose(db *bun.DB, request wspkg.CloseWrapper, llm llm.Client) {
	err := llm.Close(request.SessionID)
	if err != nil {
		utils.Log(utils.ErrorLevel).Err(err).CID(request.SessionID).BT().Send("Failed to close conversation")
	}
}
