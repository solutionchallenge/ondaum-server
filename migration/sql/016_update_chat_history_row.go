package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser016UpdateChatHistoryRow = `
UPDATE chat_histories SET role = 'model' WHERE role = 'assistant';`

var MigrationUser016UpdateChatHistoryRow = database.Migration{
	Name:  "user.016.update_chat_history_row",
	Query: sqlUser016UpdateChatHistoryRow,
}
