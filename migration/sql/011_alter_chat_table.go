package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser011AlterChatTable = `
ALTER TABLE IF EXISTS chats
ADD COLUMN chat_duration BIGINT AFTER archived_at`

var MigrationUser011AlterChatTable = database.Migration{
	Name:  "user.011.alter_chat_table",
	Query: sqlUser011AlterChatTable,
}
