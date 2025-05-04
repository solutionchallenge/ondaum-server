package user

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser004CreateUserAdditionTable = `
CREATE TABLE IF NOT EXISTS user_additions
(
    user_id BIGINT PRIMARY KEY,
	concerns JSON NOT NULL COMMENT 'array of undefined user concerns',
	emotions JSON NOT NULL COMMENT 'array of application-defined emotion enums',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`

var MigrationUser004CreateUserAdditionTable = database.Migration{
	Name:  "user.004.create_user_addition_table",
	Query: sqlUser004CreateUserAdditionTable,
}
