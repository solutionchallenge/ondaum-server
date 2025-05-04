package user

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser004CreateAdditionTable = `
CREATE TABLE IF NOT EXISTS user_additions
(
    user_id BIGINT PRIMARY KEY,
	concerns JSON NOT NULL COMMENT 'array of undefined user concerns',
	emotions JSON NOT NULL COMMENT 'array of application-defined emotion enums',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`

var MigrationUser004CreateAdditionTable = database.Migration{
	Name:  "user.004.create_addition_table",
	Query: sqlUser004CreateAdditionTable,
}
