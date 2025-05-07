package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser009CreateDiagnosisTable = `
CREATE TABLE IF NOT EXISTS diagnoses
(
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    diagnosis VARCHAR(50) NOT NULL,
    total_score BIGINT NOT NULL,
    result_score BIGINT NOT NULL,
    result_name VARCHAR(255) NOT NULL,
    result_description TEXT NOT NULL,
    result_critical BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`

var MigrationUser009CreateDiagnosisTable = database.Migration{
	Name:  "user.009.create_diagnosis_table",
	Query: sqlUser009CreateDiagnosisTable,
}
