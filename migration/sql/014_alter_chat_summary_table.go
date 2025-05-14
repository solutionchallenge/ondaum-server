package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser014AlterChatSummaryTable = `
ALTER TABLE chat_summaries
    ADD PRIMARY KEY (chat_id),
    DROP COLUMN id`

var MigrationUser014AlterChatSummaryTable = database.Migration{
	Name:  "user.014.alter_chat_summary_table",
	Query: sqlUser014AlterChatSummaryTable,
}
