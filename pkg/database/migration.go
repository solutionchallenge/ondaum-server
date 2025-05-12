package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

const createMigrationTable = `
CREATE TABLE IF NOT EXISTS migrations
(
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	migration_name VARCHAR(255) NOT NULL UNIQUE,
	migration_query TEXT NOT NULL
)`

const insertMigrationHistory = `
INSERT INTO migrations (migration_name, migration_query) VALUES (?, ?)
`

type Migration struct {
	Name  string
	Query string
}

func (m *Migration) Apply(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return utils.WrapError(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, m.Query); err != nil {
		return utils.WrapError(err, "failed to execute migration query")
	}

	if _, err := tx.ExecContext(ctx, insertMigrationHistory, m.Name, m.Query); err != nil {
		return utils.WrapError(err, "failed to insert migration history")
	}

	if err := tx.Commit(); err != nil {
		return utils.WrapError(err, "failed to commit transaction")
	}
	return nil
}

func getMigrationHistories[EXEC *sql.Tx | *sql.DB](ctx context.Context, executor EXEC) ([]MigrationHistory, error) {
	var querier func(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	switch executor := any(executor).(type) {
	case *sql.Tx:
		querier = executor.QueryContext
	case *sql.DB:
		querier = executor.QueryContext
	default:
		return nil, utils.WrapError(nil, "invalid executor type")
	}

	rows, err := querier(ctx, "SELECT migration_name, created_at FROM migrations ORDER BY created_at ASC")
	if err != nil {
		return nil, utils.WrapError(err, "failed to query migrations")
	}
	defer rows.Close()

	histories := []MigrationHistory{}
	for rows.Next() {
		var history MigrationHistory
		if err := rows.Scan(&history.Name, &history.CreatedAt); err != nil {
			return nil, utils.WrapError(err, "failed to scan migration history")
		}
		histories = append(histories, history)
	}
	if err := rows.Err(); err != nil {
		return nil, utils.WrapError(err, "error iterating migrations")
	}

	return histories, nil
}

func GetMigrationHistories(ctx context.Context, db *sql.DB) ([]MigrationHistory, error) {
	return getMigrationHistories(ctx, db)
}

func Migrate(db *sql.DB, migrations ...Migration) error {
	ctx := context.Background()
	utils.Log(utils.InfoLevel).Ctx(ctx).BT().Send("Migrating database...")

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return utils.WrapError(err, "failed to begin migration transaction")
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, createMigrationTable); err != nil {
		return utils.WrapError(err, "failed to create migration table")
	}

	histories, err := getMigrationHistories(ctx, tx)
	if err != nil {
		return utils.WrapError(err, "failed to get migration histories")
	}

	if len(histories) == 0 {
		utils.Log(utils.InfoLevel).Ctx(ctx).BT().Send("No migration history found. Applying all migrations...")
		for _, migration := range migrations {
			utils.Log(utils.InfoLevel).Ctx(ctx).BT().Send("Migrating %v...", migration.Name)
			if _, err := tx.ExecContext(ctx, migration.Query); err != nil {
				return utils.WrapError(err, "failed to execute migration query for %s", migration.Name)
			}
			if _, err := tx.ExecContext(ctx, insertMigrationHistory, migration.Name, migration.Query); err != nil {
				return utils.WrapError(err, "failed to insert migration history for %s", migration.Name)
			}
		}
	} else {
		for i, history := range histories {
			if i >= len(migrations) {
				return utils.WrapError(nil, "database has more migrations (%s) than defined in code", history.Name)
			}
			if history.Name != migrations[i].Name {
				return utils.WrapError(nil, "migration history mismatch: expected %s, got %s", migrations[i].Name, history.Name)
			}
		}

		lastAppliedIndex := len(histories)
		if lastAppliedIndex < len(migrations) {
			utils.Log(utils.InfoLevel).Ctx(ctx).BT().Send("Applying incremental migrations from %s...", migrations[lastAppliedIndex].Name)
			for i := lastAppliedIndex; i < len(migrations); i++ {
				migration := migrations[i]
				utils.Log(utils.InfoLevel).Ctx(ctx).BT().Send("Migrating %v...", migration.Name)
				if _, err := tx.ExecContext(ctx, migration.Query); err != nil {
					return utils.WrapError(err, "failed to execute migration query for %s", migration.Name)
				}
				if _, err := tx.ExecContext(ctx, insertMigrationHistory, migration.Name, migration.Query); err != nil {
					return utils.WrapError(err, "failed to insert migration history for %s", migration.Name)
				}
			}
		} else {
			utils.Log(utils.InfoLevel).Ctx(ctx).BT().Send("Database is up to date")
		}
	}

	if err := tx.Commit(); err != nil {
		return utils.WrapError(err, "failed to commit migration transaction")
	}

	return nil
}

type MigrationHistory struct {
	Name      string    `db:"migration_name"`
	CreatedAt time.Time `db:"created_at"`
}
