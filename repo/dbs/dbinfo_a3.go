package dbs

import (
	"fmt"

	"korrectkm/config"

	"github.com/mechiko/utility"
)

func NewA3(app ILogCfg, dbType, fsrarId string, cfg config.DatabaseConfiguration) (dbi *DbInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			dbi.Exists = false
		}
	}()
	// configName := "alcohelp3"
	dbi = NewDbInfo(cfg)
	if dbi.Driver = dbType; dbi.Driver == "" {
		dbi.Driver = "sqlite"
	}
	if dbi.Name == "" {
		if dbi.Name = fsrarId; dbi.Name == "" {
			if dbi.Name = fsrarId; dbi.Name == "" {
				return dbi, fmt.Errorf("%s отсутствуют имя базы данных для А3", modError)
			}
		}
	}
	if dbi.File = fsrarId; dbi.File == "" {
		return dbi, fmt.Errorf("%s отсутствуют имя базы данных для A3", modError)
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
