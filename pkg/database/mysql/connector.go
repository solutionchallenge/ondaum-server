package mysql

import (
	"fmt"

	"github.com/lithammer/shortuuid"
	"github.com/solutionchallenge/ondaum-server/pkg/database"
)

const (
	Kind = "mysql"
)

func NewConnector(config database.Config) (database.Connector, error) {
	if config.Kind != Kind {
		return database.Connector{}, fmt.Errorf("unsupported database driver: %s", config.Kind)
	}
	connstr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&interpolateParams=true&multiStatements=false",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	return database.Connector{
		Kind:             Kind,
		Identifier:       shortuuid.New(),
		ConnectionString: connstr,
	}, nil
}
