package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/dimazaicev88/ts/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewMysql(conf config.DbConfig) (*bun.DB, error) {
	urlConnect := fmt.Sprintf("%s:%s@tcp(%s:%s)",
		conf.Login,
		conf.Password,
		conf.Host,
		conf.Port,
	) // Базовое url для подключения

	if conf.AllowExecuteSQL {
		urlConnect += fmt.Sprintf("/%s?multiStatements=true", conf.DataBaseName)
	} else if len(strings.TrimSpace(conf.DataBaseName)) == 0 {
		urlConnect += "/?parseTime=true"
	} else {
		return nil, fmt.Errorf("not compitable mode work for databese")
	}

	sqlDb, err := sql.Open("mysql", urlConnect)
	if err != nil {
		return nil, err
	}
	db := bun.NewDB(sqlDb, mysqldialect.New())

	db.WithQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	if conf.SetTimezone {
		_, err = db.Exec("SET time_zone = 'Europe/Moscow'")
		if err != nil {
			return nil, fmt.Errorf("ошибка при установке часового пояса:%s", err.Error())
		}
	}

	return db, nil
}

func NewMysqlByUrl(urlConnect string) (*bun.DB, error) {
	sqlDb, err := sql.Open("mysql", urlConnect)
	if err != nil {
		return nil, err
	}
	db := bun.NewDB(sqlDb, mysqldialect.New())
	return db, nil
}
