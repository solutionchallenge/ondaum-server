package migration

import (
	"github.com/solutionchallenge/ondaum-server/migration/user"
	"github.com/solutionchallenge/ondaum-server/pkg/database"
)

var Collection = []database.Migration{
	user.MigrationUser001CreateUserTable,
	user.MigrationUser002CreateOAuthTable,
	user.MigrationUser003CreatePrivacyTable,
	user.MigrationUser004CreateAdditionTable,
	user.MigrationUser005CreateFutureJobTable,
	user.MigrationUser006CreateChatTable,
	user.MigrationUser007CreateChatHistoryTable,
	user.MigrationUser008CreateChatSummaryTable,
}
