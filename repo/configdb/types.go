package configdb

import (
	_ "embed"
	"fmt"

	"github.com/mechiko/dbscan"
	"github.com/upper/db/v4"
)

const modError = "configdb"

type DbConfig struct {
	dbSession db.Session // открытый хэндл тут
	dbInfo    *dbscan.DbInfo
	infoType  dbscan.DbInfoType
	version   int64
}

func New(info *dbscan.DbInfo) (*DbConfig, error) {
	if info == nil {
		return nil, fmt.Errorf("%s dbinfo is nil", modError)
	}
	db := &DbConfig{
		dbInfo:   info,
		infoType: dbscan.Config,
	}
	// открываем сесиию в этом методе если нет ошибки
	if err := db.Check(); err != nil {
		return nil, fmt.Errorf("%s check failed: %w", modError, err)
	}
	if db.dbSession == nil {
		return nil, fmt.Errorf("%s error after check dbsession nil", modError)
	}
	return db, nil
}

func (c *DbConfig) Close() error {
	if c.dbSession == nil {
		return nil
	}
	err := c.dbSession.Close()
	c.dbSession = nil
	return err
}

func (c *DbConfig) Sess() db.Session {
	return c.dbSession
}

func (c *DbConfig) Version() int64 {
	return c.version
}

func (c *DbConfig) Info() dbscan.DbInfo {
	if c.dbInfo == nil {
		return dbscan.DbInfo{}
	}
	return *c.dbInfo
}

func (c *DbConfig) InfoType() dbscan.DbInfoType {
	return c.infoType
}
