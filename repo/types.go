package repo

import (
	"context"
	"fmt"

	"korrectkm/config"
	"korrectkm/repo/a3"
	"korrectkm/repo/configdb"
	"korrectkm/repo/dbs"
	"korrectkm/repo/selfdb"
	"korrectkm/repo/znakdb"

	"github.com/mechiko/utility"

	"go.uber.org/zap"
)

const modError = "pkg:repo"

type ILogCfg interface {
	Config() config.IConfig
	Logger() *zap.SugaredLogger
}

type Repo interface {
	Run(context.Context) error
	Shutdown()
	Self() *selfdb.DbSelf
	ConfigDB() *configdb.DbConfig
}

var Version int64

type Repository struct {
	ILogCfg
	dbs     *dbs.Dbs
	fsrarId string
	errors  []string
}

var _ Repo = (*Repository)(nil)

// dbPath для своей БД
// func New(logcfg ILogCfg, dbPath string) (rp *Repository, err error) {
func New(logcfg ILogCfg, dbPath string) (rp *Repository) {
	defer func() {
		// ошибки дописываем в массив
		// реально ошибка не возвращается ни когда TODO
		if r := recover(); r != nil {
			errStr := fmt.Sprintf("доступ к бд %v", r)
			rp.errors = append(rp.errors, errStr)
		}
	}()

	rp = &Repository{
		ILogCfg: logcfg,
		errors:  make([]string, 0),
		dbs:     dbs.New(logcfg, "config.db", dbPath),
	}
	// создаем объект описателей БД
	// имя БД конфигурации и с признаком сканирования других БД
	// rp.dbs = dbs.New(rp, "config.db", dbPath)
	if len(rp.dbs.Errors()) > 0 {
		rp.errors = append(rp.errors, rp.dbs.Errors()...)
		// return
	}
	cfg := logcfg.Config().Configuration()
	if utility.StructHasField(cfg, "Alcohelp3") {
		if !rp.A3DBPing() {
			rp.Logger().Infof("%s отсутствует БД А3", modError)
			rp.dbs.A3().Exists = false
		}
		switch rp.dbs.A3().Driver {
		case "sqlite":
			dbname := utility.FileNameWithoutExtension(rp.dbs.A3().File)
			if rp.Config().Configuration().Application.Fsrarid != dbname {
				rp.Config().SetInConfig("application.fsrarid", dbname, true)
			}
		case "mssql":
			if rp.Config().Configuration().Application.Fsrarid != rp.dbs.A3().Name {
				rp.Config().SetInConfig("application.fsrarid", rp.dbs.A3().Name, true)
			}
		}

	}
	if utility.StructHasField(cfg, "Config") {
		if !rp.ConfigDBPing() {
			// когда бд config.db нет что делаем
			rp.Logger().Infof("%s отсутствует БД config", modError)
			rp.dbs.ConfigInfo().Exists = false
		}
	}
	if utility.StructHasField(cfg, "TrueZnak") {
		if !rp.ZnakDBPing() {
			// когда бд 4z db нет что делаем
			rp.Logger().Infof("%s отсутствует БД 4z", modError)
			rp.dbs.Znak().Exists = false
		}
	}

	if utility.StructHasField(cfg, "SelfDB") {
		// инициализируем или проверяем БД self
		if err := rp.prepareSelf(); err != nil {
			rp.Logger().Infof("%s подготовка %s.db %s", modError, config.Name, err.Error())
			rp.AddError(err.Error())
		}
		if !rp.SelfDBPing() {
			// когда бд 4z db нет что делаем
			rp.Logger().Infof("%s отсутствует БД %s", modError, rp.dbs.Self().Name)
			rp.dbs.Self().Exists = false
			rp.AddError("selfdb отсутствует")
		}
	}

	// это заглушка пока не работает просто ни чего не делает
	if err := rp.dbs.SaveConfig(); err != nil {
		rp.AddError(err.Error())
		// return nil, fmt.Errorf("%s %w", modError, err)
	}
	return rp
}

func (r *Repository) Dbs() *dbs.Dbs {
	return r.dbs
}

func (r *Repository) FsrarID() string {
	return r.fsrarId
}

func (r *Repository) Self() *selfdb.DbSelf {
	return selfdb.New(r, r.dbs.Self())
}

func (r *Repository) ConfigDB() *configdb.DbConfig {
	return configdb.New(r, r.dbs.ConfigInfo())
}

func (r *Repository) ZnakDB() *znakdb.DbZnak {
	return znakdb.New(r, r.dbs.Znak())
}

func (r *Repository) A3DB() *a3.DbA3 {
	return a3.New(r, r.dbs.A3())
}

func (r *Repository) Errors() []string {
	return r.errors
}

func (r *Repository) AddError(e string) {
	r.errors = append(r.errors, e)
}
