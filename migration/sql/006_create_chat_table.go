package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser006CreateChatTable = `
CREATE TABLE IF NOT EXISTS chats
(
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    session_id VARCHAR(50) NOT NULL,
    started_date DATETIME NOT NULL,
    user_timezone VARCHAR(50) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    finished_at DATETIME,
    archived_at DATETIME,
    UNIQUE KEY unique_session_user (session_id, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`

var MigrationUser006CreateChatTable = database.Migration{
	Name:  "user.006.create_chat_table",
	Query: sqlUser006CreateChatTable,
}
