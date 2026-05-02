package fixtures

import (
	"context"
	"os"

	"github.com/dimazaicev88/ts/internal/storage"

	"github.com/uptrace/bun/extra/bundebug"

	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/uptrace/bun"
)

func setLogLevel(db *bun.DB) error {
	bunDebugKey := "BUNDEBUG"
	err := os.Setenv(bunDebugKey, "2")
	if err != nil {
		return err
	}
	db.WithQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv(bunDebugKey),
	))
	return nil
}

func SetupTestContainerMysql(
	ctx context.Context,
	testContainer *mysql.MySQLContainer,
	pathToMigrations, dataBaseName string,
) (*bun.DB, error) {
	urlConnect, err := testContainer.ConnectionString(ctx, "multiStatements=true")

	if err != nil {
		return nil, err
	}
	db, err := storage.NewMysqlByUrl(urlConnect)

	if err != nil {
		return nil, err
	}
	err = storage.MigrateDbByUrlConnect(urlConnect, pathToMigrations, dataBaseName)

	if err != nil {
		return nil, err
	}

	err = setLogLevel(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetupByUrlConnect(
	ctx context.Context,
	urlConnect,
	pathToMigrations, dataBaseName string,
) (*bun.DB, error) {
	db, err := storage.NewMysqlByUrl(urlConnect)

	if err != nil {
		return nil, err
	}
	err = storage.MigrateDbByUrlConnect(urlConnect, pathToMigrations, dataBaseName)

	if err != nil {
		return nil, err
	}

	err = setLogLevel(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
