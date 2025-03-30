package memdb

import (
	"context"
	"fmt"
	"net"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/lithammer/shortuuid/v3"
	"github.com/phayes/freeport"
	"github.com/solutionchallenge/ondaum-server/pkg/database"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

const (
	Kind = "memdb"
)

func checkConnectability(ctx context.Context, retries int, port int) error {
	return utils.Retry(ctx, retries, func() error {
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			return err
		}
		defer conn.Close()
		return nil
	}, func(err error) bool {
		return err != nil
	})
}

func initializeServer(dbname string) (*memory.DbProvider, *sqle.Engine, error) {
	db := memory.NewDatabase(dbname)
	db.BaseDatabase.EnablePrimaryKeyIndexes()
	provider := memory.NewDBProvider(db)
	engine := sqle.NewDefault(provider)

	session := memory.NewSession(sql.NewBaseSession(), provider)
	ctx := sql.NewContext(context.Background(), sql.WithSession(session))
	ctx.SetCurrentDatabase(dbname)

	return provider, engine, nil
}

func startServer(port int, dbname string) error {
	provider, engine, err := initializeServer(dbname)
	if err != nil {
		return fmt.Errorf("failed to initialize memdb server: %w", err)
	}

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("127.0.0.1:%d", port),
	}

	srv, err := server.NewServer(config, engine, memory.NewSessionBuilder(provider), nil)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	if err = srv.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func NewConnector(config database.Config) (database.Connector, error) {
	if config.Kind != Kind {
		return database.Connector{}, fmt.Errorf("unsupported database driver: %s", config.Kind)
	}

	port := config.Port

	if port == 0 {
		free, err := freeport.GetFreePort()
		if err != nil {
			return database.Connector{}, fmt.Errorf("failed to get free port: %w", err)
		}
		port = free
	}

	go func() {
		if err := startServer(port, config.Database); err != nil {
			panic(fmt.Sprintf("failed to start memdb server: %v", err))
		}
	}()

	if err := checkConnectability(context.Background(), 5, port); err != nil {
		return database.Connector{}, fmt.Errorf("failed to connect to memdb: %w", err)
	}

	return database.Connector{
		Kind:       "mysql",
		Identifier: shortuuid.New(),
		ConnectionString: fmt.Sprintf(
			"root:root@tcp(localhost:%d)/%s?parseTime=true&multiStatements=true",
			port, config.Database,
		),
	}, nil
}
