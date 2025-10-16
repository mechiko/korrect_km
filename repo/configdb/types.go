package configdb

import (
	"database/sql"
	_ "embed"
	"fmt"
	"korrectkm/config"
	"korrectkm/repo/dbs"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mssql"
	"github.com/upper/db/v4/adapter/sqlite"
	"go.uber.org/zap"
)

type ILogCfg interface {
	Config() config.IConfig
	Logger() *zap.SugaredLogger
}

type DbConfig struct {
	ILogCfg
	dbSession db.Session // открытый хэндл тут
}

const modError = "repo:configdb"

func New(logcfg ILogCfg, a *dbs.DbInfo) *DbConfig {
	db := &DbConfig{
		ILogCfg: logcfg,
	}
	switch a.Driver {
	case "mssql":
		uri := mssql.ConnectionURL{
			User:     "",
			Password: "",
			Host:     "localhost:1433",
			Database: a.Name,
			Options: map[string]string{
				"encrypt": "disable",
			},
		}
		dbSess, err := mssql.Open(uri)
		if err != nil {
			panic(fmt.Sprintf("mssql config.db %s", err.Error()))
		}
		db.dbSession = dbSess
		return db
	case "sqlite":
		uri := sqlite.ConnectionURL{
			Database: a.File,
			Options: map[string]string{
				"mode": "rw",
				// "_journal_mode": "WAL",
			},
		}
		dbSess, err := sqlite.Open(uri)
		if err != nil {
			panic(fmt.Sprintf("sqlite config.db %s", err.Error()))
		}
		db.dbSession = dbSess
		return db
	}
	panic("config.db не указан драйвер")
}

func (c *DbConfig) Close() (err error) {
	return c.dbSession.Close()
}

func (c *DbConfig) DB() *sql.DB {
	return c.dbSession.Driver().(*sql.DB)
}

func (c *DbConfig) Sess() db.Session {
	return c.dbSession
}

func (c *DbConfig) Ping() (err error) {
	sess := c.dbSession
	return sess.Ping()
}
