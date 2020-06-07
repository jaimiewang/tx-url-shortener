package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v2"
	"strings"
	"time"
	"tx-url-shortener/config"
)

var DbMap *gorp.DbMap

func buildMySQLDataSourceName(user string, pass string, hostname string, port uint16, name string) string {
	builder := strings.Builder{}
	builder.WriteString(user)
	builder.WriteRune(':')
	builder.WriteString(pass)
	builder.WriteString("@tcp(")
	builder.WriteString(hostname)
	builder.WriteRune(':')
	builder.WriteString(string(port))
	builder.WriteString(")/")
	builder.WriteString(name)

	return builder.String()
}

func InitDatabase() error {
	var driverName string
	var dataSourceName string
	var dialect gorp.Dialect

	switch dbEngine := strings.ToLower(config.Config.Database.Engine); dbEngine {
	case "sqlite3", "sqlite":
		driverName = "sqlite3"
		dataSourceName = config.Config.Database.Name
		dialect = gorp.SqliteDialect{}
	case "mysql":
		driverName = "mysql"
		dataSourceName = buildMySQLDataSourceName(
			config.Config.Database.User,
			config.Config.Database.Password,
			config.Config.Database.Hostname,
			config.Config.Database.Port,
			config.Config.Database.Name,
		)
		dialect = gorp.MySQLDialect{}
	default:
		return fmt.Errorf("not supported database engine: %s", dbEngine)
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	DbMap = &gorp.DbMap{Db: db, Dialect: dialect}
	return nil
}
