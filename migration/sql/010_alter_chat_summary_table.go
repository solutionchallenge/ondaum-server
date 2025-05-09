package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser010AlterChatSummaryTable = `
ALTER TABLE IF EXISTS chat_summaries
ADD COLUMN recommendations JSON AFTER emotions`

var MigrationUser010AlterChatSummaryTable = database.Migration{
	Name:  "user.010.alter_chat_summary_table",
	Query: sqlUser010AlterChatSummaryTable,
}
