package migration

import (
	"github.com/solutionchallenge/ondaum-server/migration/sql"
	"github.com/solutionchallenge/ondaum-server/pkg/database"
)

var Collection = []database.Migration{
	sql.MigrationUser001CreateUserTable,
	sql.MigrationUser002CreateOAuthTable,
	sql.MigrationUser003CreatePrivacyTable,
	sql.MigrationUser004CreateAdditionTable,
	sql.MigrationUser005CreateFutureJobTable,
	sql.MigrationUser006CreateChatTable,
	sql.MigrationUser007CreateChatHistoryTable,
	sql.MigrationUser008CreateChatSummaryTable,
	sql.MigrationUser009CreateDiagnosisTable,
	sql.MigrationUser010AlterChatSummaryTable,
	sql.MigrationUser011AlterChatTable,
	sql.MigrationUser012AlterChatSummaryTable,
	sql.MigrationUser013AlterDiagnosisTable,
	sql.MigrationUser014AlterDiagnosisTable,
}
