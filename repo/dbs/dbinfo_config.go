package dbs

import (
	"fmt"
	"korrectkm/config"

	"github.com/mechiko/utility"
)

const defaultConfigFile = "config.db"
const defaultConfigDriver = "sqlite"

// база config.db всегда sqlite3 пока
func NewConfig(app ILogCfg, cfg config.DatabaseConfiguration) (dbi *DbInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	driver := cfg.Driver
	if driver == "" {
		driver = defaultConfigDriver
	}
	dbi = &DbInfo{
		File:   defaultConfigFile,
		Driver: driver,
		Host:   cfg.Host,
		User:   cfg.User,
		Exists: false,
	}
	dbi.Exists = utility.PathOrFileExists(dbi.File)
	return dbi, nil
}
