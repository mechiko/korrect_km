package dbs

import (
	"fmt"
	"korrectkm/config"

	"github.com/mechiko/utility"
)

func New4z(app ILogCfg, dbType, name string, cfg config.DatabaseConfiguration) (dbi *DbInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			dbi.Exists = false
		}
	}()
	dbi = NewDbInfo(cfg)
	if dbi.Driver = dbType; dbi.Driver == "" {
		dbi.Driver = "sqlite"
	}
	if dbi.File = name; dbi.File == "" {
		return dbi, fmt.Errorf("%s отсутствуют имя базы данных для 4Z", modError)
	}
	dbi.File = dbi.File + ".db"
	if dbi.Driver == "sqlite" {
		dbi.Exists = utility.PathOrFileExists(dbi.File)
	}
	if dbi.Driver == "mssql" {
		if cfg.Host == "" {
			dbi.Host = "localhost:1433"
		}
		dbi.Host = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		dbi.User = cfg.User
		dbi.Pass = cfg.Pass
	}
	return dbi, nil
}
