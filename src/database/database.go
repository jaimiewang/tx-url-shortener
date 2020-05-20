package database

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v2"
	"strings"
	"tx-url-shortener/config"
	"tx-url-shortener/model"
)

var DbMap *gorp.DbMap

func InitDatabase() error {
	var driverName string
	var dataSourceName string
	var dialect gorp.Dialect

	switch strings.ToLower(config.Conf.Database.Engine) {
	case "sqlite3":
		driverName = "sqlite3"
		dataSourceName = config.Conf.Database.Name
		dialect = gorp.SqliteDialect{}
	default:
		return errors.New("invalid database engine")
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	DbMap = &gorp.DbMap{Db: db, Dialect: dialect}
	DbMap.AddTableWithName(model.ShortURL{}, "urls")

	err = DbMap.CreateTablesIfNotExists()
	if err != nil {
		return err
	}

	return nil
}
