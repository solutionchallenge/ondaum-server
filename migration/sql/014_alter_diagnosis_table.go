package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser014AlterDiagnosisTable = `
ALTER TABLE diagnoses
MODIFY COLUMN sub_id VARCHAR(255) NOT NULL AFTER id,
ADD INDEX idx_diagnoses_user_sub_id (user_id, sub_id)`

var MigrationUser014AlterDiagnosisTable = database.Migration{
	Name:  "user.014.alter_diagnosis_table",
	Query: sqlUser014AlterDiagnosisTable,
}
