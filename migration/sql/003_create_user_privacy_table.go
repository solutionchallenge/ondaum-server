package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser003CreatePrivacyTable = `
CREATE TABLE IF NOT EXISTS user_privacies
(
    user_id BIGINT PRIMARY KEY,
	gender ENUM('male', 'female', 'other') NOT NULL,
	birthday DATE NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`

var MigrationUser003CreatePrivacyTable = database.Migration{
	Name:  "user.003.create_privacy_table",
	Query: sqlUser003CreatePrivacyTable,
}
