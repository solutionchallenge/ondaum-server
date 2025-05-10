package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser012AlterChatSummaryTable = `
ALTER TABLE chat_summaries
ADD COLUMN positive_score DOUBLE AFTER recommendations,
ADD COLUMN negative_score DOUBLE AFTER positive_score,
ADD COLUMN neutral_score DOUBLE AFTER negative_score,
ADD COLUMN main_topic JSON AFTER neutral_score`

var MigrationUser012AlterChatSummaryTable = database.Migration{
	Name:  "user.012.alter_chat_summary_table",
	Query: sqlUser012AlterChatSummaryTable,
}
