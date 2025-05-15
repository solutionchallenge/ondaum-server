package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser015AlterChatTable = `
ALTER TABLE chats
DROP COLUMN started_date,
ADD INDEX idx_user_session (user_id, session_id),
ADD INDEX idx_user_created_archived (user_id, created_at, archived_at);`

var MigrationUser015AlterChatTable = database.Migration{
	Name:  "user.015.alter_chat_table",
	Query: sqlUser015AlterChatTable,
}
