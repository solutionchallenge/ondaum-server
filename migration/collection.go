package migration

import (
	"github.com/solutionchallenge/ondaum-server/migration/user"
	"github.com/solutionchallenge/ondaum-server/pkg/database"
)

var Collection = []database.Migration{
	user.MigrationUser001CreateUserTable,
	user.MigrationUser002CreateOAuthTable,
}
