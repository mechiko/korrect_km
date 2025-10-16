package dbs

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/sqlite"
)

type Parameters struct {
	Name  string `db:"name"`
	Value string `db:"value"`
}

func (p *Parameters) Store(sess db.Session) db.Store {
	return sess.Collection("parameters")
}

func (dd *Dbs) fromConfig(key string) (out string) {
	var dbs db.Session
	var err error
	defer func() {
		if r := recover(); r != nil {
			out = ""
		}
	}()

	if dd.config == nil {
		return ""
	}
	info := dd.config
	switch dd.config.Driver {
	case "sqlite":
		uri := sqlite.ConnectionURL{
			Database: info.File,
			Options: map[string]string{
				"mode": "rw",
			},
		}
		dbs, err = sqlite.Open(uri)
		if err != nil {
			dd.Logger().Errorf("dbs:fromconfig %s", err.Error())
			return ""
		}
		defer func() {
			if err = dbs.Close(); err != nil {
				dd.Logger().Errorf("dbs:fromconfig %s", err.Error())
			}
		}()
	}
	param := &Parameters{}
	if err = dbs.Get(param, db.Cond{"name": key}); err != nil {
		dd.Logger().Errorf("dbs:fromconfig %s", err.Error())
	}
	return param.Value
}
