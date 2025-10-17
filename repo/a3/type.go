package a3

import (
	"fmt"

	"github.com/mechiko/dbscan"
	"github.com/upper/db/v4"
)

const modError = "a3db"

type DbA3 struct {
	dbSession db.Session // открытый хэндл тут
	dbInfo    *dbscan.DbInfo
	infoType  dbscan.DbInfoType
	version   int64
}

func New(info *dbscan.DbInfo) (*DbA3, error) {
	if info == nil {
		return nil, fmt.Errorf("%s dbinfo is nil", modError)
	}
	db := &DbA3{
		dbInfo:   info,
		infoType: dbscan.A3,
	}
	// открываем сесиию в этом методе если нет ошибки
	if err := db.Check(); err != nil {
		return nil, fmt.Errorf("%s error check %w", modError, err)
	}
	if db.dbSession == nil {
		return nil, fmt.Errorf("%s error after check dbsession nil", modError)
	}
	return db, nil
}

func (c *DbA3) Close() (err error) {
	if c.dbSession == nil {
		return nil
	}
	err = c.dbSession.Close()
	c.dbSession = nil
	return err
}

func (c *DbA3) Sess() db.Session {
	return c.dbSession
}

func (c *DbA3) Version() int64 {
	return c.version
}

func (c *DbA3) Info() dbscan.DbInfo {
	if c.dbInfo == nil {
		return dbscan.DbInfo{}
	}
	return *c.dbInfo
}

func (c *DbA3) InfoType() dbscan.DbInfoType {
	return c.infoType
}
