package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser008CreateChatSummaryTable = `
CREATE TABLE IF NOT EXISTS chat_summaries
(
    chat_id BIGINT PRIMARY KEY,
    title VARCHAR(255),
    text TEXT,
    keywords JSON,
    emotions JSON,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
)`

var MigrationUser008CreateChatSummaryTable = database.Migration{
	Name:  "user.008.create_chat_summary_table",
	Query: sqlUser008CreateChatSummaryTable,
}
