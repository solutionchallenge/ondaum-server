package user

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser007CreateChatHistoryTable = `
CREATE TABLE IF NOT EXISTS chat_histories
(
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    metadata JSON,
    message_id VARCHAR(50) NOT NULL,
    inserted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    chat_id BIGINT NOT NULL,
    UNIQUE KEY unique_chat_message (chat_id, message_id),
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
)`

var MigrationUser007CreateChatHistoryTable = database.Migration{
	Name:  "user.007.create_chat_history_table",
	Query: sqlUser007CreateChatHistoryTable,
}
