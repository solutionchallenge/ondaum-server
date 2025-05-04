package user

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser003CreateUserPrivaciesTable = `
CREATE TABLE IF NOT EXISTS user_privacies
(
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
	gender ENUM('male', 'female', 'other') NOT NULL,
	birthday DATE NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`

var MigrationUser003CreateUserPrivaciesTable = database.Migration{
	Name:  "user.003.create_user_privacies_table",
	Query: sqlUser003CreateUserPrivaciesTable,
}
