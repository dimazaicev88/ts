package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dimazaicev88/ts/internal/storage/models"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

func TruncateAllTables(ctx context.Context, db *bun.DB) {
	tables := []any{
		new(models.Task),
		new(models.Worker),
	}
	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		logrus.Fatal(err.Error())
	}
	for _, table := range tables {
		_, err := db.NewTruncateTable().Model(table).Exec(ctx)
		if err != nil {
			logrus.Fatal(err.Error())
		}
	}
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		logrus.Fatal(err.Error())
	}
}

func TruncateTable(ctx context.Context, db *bun.DB, table string) error {
	_, err := db.NewRaw(fmt.Sprintf("SET FOREIGN_KEY_CHECKS = 0; TRUNCATE TABLE `%s`; SET FOREIGN_KEY_CHECKS = 1;", table)).Exec(ctx)
	return err
}

func IsErrNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsNoErrNoRows(err error) bool {
	return !errors.Is(err, sql.ErrNoRows)
}
