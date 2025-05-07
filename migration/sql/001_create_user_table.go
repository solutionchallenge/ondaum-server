package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser001CreateUserTable = `
CREATE TABLE IF NOT EXISTS users
(
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    attempts BIGINT NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`

var MigrationUser001CreateUserTable = database.Migration{
	Name:  "user.001.create_user_table",
	Query: sqlUser001CreateUserTable,
}
