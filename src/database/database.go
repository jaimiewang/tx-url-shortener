package database

import (
	"database/sql"
	"errors"
	"github.com/coocood/freecache"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v2"
	"strings"
	"tx-url-shortener/config"
)

var DbMap *gorp.DbMap
var Cache *freecache.Cache

func InitDatabase() error {
	var driverName string
	var dataSourceName string
	var dialect gorp.Dialect

	switch strings.ToLower(config.Config.Database.Engine) {
	case "sqlite3":
		driverName = "sqlite3"
		dataSourceName = config.Config.Database.Name
		dialect = gorp.SqliteDialect{}
	default:
		return errors.New("invalid database engine")
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	DbMap = &gorp.DbMap{Db: db, Dialect: dialect}
	Cache = freecache.NewCache(config.Config.CacheSize)
	return nil
}
