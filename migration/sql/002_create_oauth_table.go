package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser002CreateOAuthTable = `
CREATE TABLE IF NOT EXISTS user_oauths
(
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    provider VARCHAR(50) NOT NULL,
    provider_code VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_provider (user_id, provider),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`

var MigrationUser002CreateOAuthTable = database.Migration{
	Name:  "user.002.create_oauth_table",
	Query: sqlUser002CreateOAuthTable,
}
