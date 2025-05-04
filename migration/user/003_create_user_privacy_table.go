package user

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser003CreateUserPrivacyTable = `
CREATE TABLE IF NOT EXISTS user_privacies
(
    user_id BIGINT PRIMARY KEY,
	gender ENUM('male', 'female', 'other') NOT NULL,
	birthday DATE NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`

var MigrationUser003CreateUserPrivacyTable = database.Migration{
	Name:  "user.003.create_user_privacy_table",
	Query: sqlUser003CreateUserPrivacyTable,
}
