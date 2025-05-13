package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser013AlterDiagnosisTable = `
ALTER TABLE diagnoses
ADD COLUMN sub_id VARCHAR(255) NOT NULL DEFAULT (UUID()) AFTER id,
ADD INDEX idx_diagnoses_user_sub_id (user_id, sub_id)`

var MigrationUser013AlterDiagnosisTable = database.Migration{
	Name:  "user.013.alter_diagnosis_table",
	Query: sqlUser013AlterDiagnosisTable,
}
