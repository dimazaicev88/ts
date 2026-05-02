package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"path"

	"github.com/dimazaicev88/ts/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func MigrateDb(config config.DbConfig, pathToMigrations string) error {
	urlConnect := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		config.Login,
		config.Password,
		config.Host,
		config.Port,
		config.DataBaseName,
	)

	db, err := sql.Open("mysql", urlConnect)
	if err != nil {
		return err
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	if len(pathToMigrations) == 0 {
		return errors.New("pathToMigrations is empty")
	}

	m, err := migrate.NewWithDatabaseInstance(
		path.Join("file://", pathToMigrations),
		config.DataBaseName,
		driver,
	)

	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func MigrateDbByUrlConnect(urlConnect, pathToMigrations, dbName string) error {
	db, err := sql.Open("mysql", urlConnect)
	if err != nil {
		logrus.Panic(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	if len(pathToMigrations) == 0 {
		logrus.Panic("PATH_TO_DB_MIGRATIONS is empty")
	}

	m, err := migrate.NewWithDatabaseInstance(
		path.Join("file://", pathToMigrations),
		dbName,
		driver,
	)

	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
