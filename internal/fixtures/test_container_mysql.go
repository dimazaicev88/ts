package fixtures

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewMysqlTestContainer(ctx context.Context, req *require.Assertions, dbName string) *mysql.MySQLContainer {
	// Общие настройки контейнера
	options := []testcontainers.ContainerCustomizer{
		mysql.WithDatabase(dbName),
		mysql.WithUsername("root"),
		mysql.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306").WithOccurrence(1).WithStartupTimeout(240 * time.Second),
		),
	}
	mysqlContainer, err := mysql.Run(ctx, "mysql:8.4.9", options...)
	req.Nil(err)
	return mysqlContainer
}
