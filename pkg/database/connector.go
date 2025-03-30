package database

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	sqlhook "github.com/qustavo/sqlhooks/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/schema"
)

type Connector struct {
	Kind             string
	Identifier       string
	ConnectionString string
}

func (connector Connector) ToSqlDB(hook ...sqlhook.Hooks) (*sql.DB, error) {
	boundDriver := string(connector.Kind)
	if len(hook) > 0 && hook[0] != nil {
		boundDriver = fmt.Sprintf("gosql-%v-%s", boundDriver, connector.Identifier)
		sql.Register(boundDriver, sqlhook.Wrap(&mysql.MySQLDriver{}, hook[0]))
	}
	sqlDB, err := sql.Open(boundDriver, connector.ConnectionString)
	if err != nil {
		return nil, err
	}
	return sqlDB, nil
}

func (connector Connector) ToSqlxDB(hook ...sqlhook.Hooks) (*sqlx.DB, error) {
	boundDriver := string(connector.Kind)
	if len(hook) > 0 && hook[0] != nil {
		boundDriver = fmt.Sprintf("sqlx-%v-%s", boundDriver, connector.Identifier)
		sql.Register(boundDriver, sqlhook.Wrap(&mysql.MySQLDriver{}, hook[0]))
	}
	sqlxDB, err := sqlx.Open(boundDriver, connector.ConnectionString)
	if err != nil {
		return nil, err
	}
	return sqlxDB, nil
}

func (connector Connector) ToBunDB(hook ...sqlhook.Hooks) (*bun.DB, error) {
	var dialect schema.Dialect
	switch connector.Kind {
	case "mysql":
		dialect = mysqldialect.New()
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", connector.Kind)
	}
	boundDriver := string(connector.Kind)
	if len(hook) > 0 && hook[0] != nil {
		boundDriver = fmt.Sprintf("bun-%v-%s", boundDriver, connector.Identifier)
		sql.Register(boundDriver, sqlhook.Wrap(&mysql.MySQLDriver{}, hook[0]))
	}
	db, err := sql.Open(boundDriver, connector.ConnectionString)
	if err != nil {
		return nil, err
	}
	return bun.NewDB(db, dialect), nil
}
