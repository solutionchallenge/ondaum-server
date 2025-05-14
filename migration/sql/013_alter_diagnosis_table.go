package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser013AlterDiagnosisTable = `
ALTER TABLE diagnoses
ADD COLUMN sub_id VARCHAR(255) NOT NULL DEFAULT ('') AFTER id;
`

var MigrationUser013AlterDiagnosisTable = database.Migration{
	Name:  "user.013.alter_diagnosis_table",
	Query: sqlUser013AlterDiagnosisTable,
}
