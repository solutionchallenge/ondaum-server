package sql

import "github.com/solutionchallenge/ondaum-server/pkg/database"

const sqlUser014AlterChatSummaryTable = `
-- Create temporary table with the most recent records
CREATE TEMPORARY TABLE temp_chat_summaries AS
SELECT DISTINCT chat_id, title, text, keywords, emotions, created_at, updated_at
FROM chat_summaries cs1
WHERE created_at = (
    SELECT MAX(created_at)
    FROM chat_summaries cs2
    WHERE cs2.chat_id = cs1.chat_id
);

-- Clear the original table
TRUNCATE TABLE chat_summaries;

-- Insert the deduplicated records back
INSERT INTO chat_summaries (chat_id, title, text, keywords, emotions, created_at, updated_at)
SELECT chat_id, title, text, keywords, emotions, created_at, updated_at
FROM temp_chat_summaries;

-- Drop temporary table
DROP TEMPORARY TABLE temp_chat_summaries;

-- Add primary key and drop id column
ALTER TABLE chat_summaries
    ADD PRIMARY KEY (chat_id),
    DROP COLUMN id`

var MigrationUser014AlterChatSummaryTable = database.Migration{
	Name:  "user.014.alter_chat_summary_table",
	Query: sqlUser014AlterChatSummaryTable,
}
