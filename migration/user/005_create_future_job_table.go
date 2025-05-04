package user

import (
	"github.com/solutionchallenge/ondaum-server/pkg/database"
	dbfuture "github.com/solutionchallenge/ondaum-server/pkg/future/database"
)

var MigrationUser005CreateFutureJobTable = database.Migration{
	Name:  "user.005.create_future_job_table",
	Query: dbfuture.FutureJobTableCreationSQL,
}
