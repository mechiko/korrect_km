package dbs

import (
	"fmt"
	"korrectkm/config"
	"path/filepath"

	"github.com/mechiko/utility"
)

// const defaultSelfDBName = "self"
const defaultSelfDBDriver = "sqlite"

func NewSelf(app ILogCfg, dbname, dbPath string, cfg config.DatabaseConfiguration) (dbi *DbInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			dbi.Exists = false
		}
	}()
	dbi = NewDbInfo(cfg)
	if dbi.Driver == "" {
		dbi.Driver = defaultSelfDBDriver
	}
	if dbi.Name == "" {
		if dbi.Name = config.Name; dbi.Name == "" {
			return dbi, fmt.Errorf("%s отсутствуют имя базы данных для Self", modError)
		}
	}
	file := filepath.Join(dbPath, fmt.Sprintf("%s.db", dbi.Name))
	if dbi.File == "" {
		if dbi.File = file; dbi.File == "" {
			return nil, fmt.Errorf("%s отсутствуют имя базы данных для Self", modError)
		}
	}
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
