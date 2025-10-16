package dbs

import (
	"fmt"
	"korrectkm/config"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mechiko/utility"

	"go.uber.org/zap"
)

const modError = "repo:dbs"

type IConfig interface {
	Set(key string, value interface{}, save ...bool) error
}

type ILogCfg interface {
	Config() config.IConfig
	Logger() *zap.SugaredLogger
}

type Dbs struct {
	ILogCfg
	// defaultDriver  string
	self           *DbInfo
	a3             *DbInfo
	znak           *DbInfo
	config         *DbInfo
	configFileName string // config.db алкохелпа
	errors         []string
}

// dbPath для своей БД
func New(logcfg ILogCfg, configFileName string, dbPath string) (d *Dbs) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			errStr := fmt.Sprintf("%s %v", modError, r)
			d.errors = append(d.errors, errStr)
		}
	}()

	d = &Dbs{
		ILogCfg:        logcfg,
		configFileName: configFileName,
		errors:         make([]string, 0),
		// defaultDriver:  logcfg.Config().Configuration().Application.DbType,
	}
	cfg := logcfg.Config().Configuration()
	fsrarid := findA3Name()
	file4z := ""
	dbType := "sqlite"
	// по возможности получаем имена в config.db
	if utility.StructHasField(cfg, "Config") {
		if d.config, err = NewConfig(d, d.infoDatabaseByKey("config")); err != nil {
			// ошибки просто логируем нам нужен полный реестр БД
			d.Logger().Errorf("%s %w", modError, err)
			// d.AddError(err.Error())
		} else {
			if d.config.Exists {
				// если есть база config.db то пытаемся найти настройки
				file4z = d.fromConfig("oms_id")
				dbType = strings.ToLower(d.fromConfig("db_type"))
			}
		}
	}
	// если структура конфига по config\toml.go
	if utility.StructHasField(cfg, "Alcohelp3") {
		if d.a3, err = NewA3(d, dbType, fsrarid, d.infoDatabaseByKey("alcohelp3")); err != nil {
			d.Logger().Errorf("%s %s", modError, err.Error())
			// d.AddError(err.Error())
		}
	}
	if utility.StructHasField(cfg, "TrueZnak") {
		if d.znak, err = New4z(d, dbType, file4z, d.infoDatabaseByKey("trueznak")); err != nil {
			d.Logger().Errorf("%s %s", modError, err.Error())
			// d.AddError(err.Error())
		}
	}
	if utility.StructHasField(cfg, "SelfDB") {
		if d.self, err = NewSelf(d, "", dbPath, d.infoDatabaseByKey("selfdb")); err != nil {
			d.Logger().Errorf("%s %s", modError, err.Error())
			// d.AddError(err.Error())
		}
	}
	return d
}

func (d *Dbs) Self() *DbInfo {
	return d.self
}

func (d *Dbs) Znak() *DbInfo {
	return d.znak
}

func (d *Dbs) A3() *DbInfo {
	return d.a3
}

func (d *Dbs) ConfigInfo() *DbInfo {
	return d.config
}

// ни чего не делаем TODO
func (d *Dbs) SaveConfig() (err error) {
	return nil
}

func findA3DbName() string {
	re, err := regexp.Compile(`^0[0-9]{11}\.db$`)
	if err != nil {
		return ""
	}
	if files, err := utility.FilteredSearchOfDirectoryTree(re, ""); err != nil {
		return ""
	} else {
		if len(files) == 0 {
			return ""
		}
		return files[0]
	}
}

func findA3Name() string {
	findName := findA3DbName()
	if findName == "" {
		return ""
	}
	_, file := filepath.Split(findName)
	before := file[0 : len(file)-len(filepath.Ext(file))]
	// before, _ := strings.CutSuffix(file, filepath.Ext(file))
	return before
}

func (d *Dbs) infoDatabaseByKey(key string) config.DatabaseConfiguration {
	return config.DatabaseConfiguration{
		Driver: d.Config().GetKeyString(fmt.Sprintf("%s.driver", key)),
		File:   d.Config().GetKeyString(fmt.Sprintf("%s.file", key)),
		DbName: d.Config().GetKeyString(fmt.Sprintf("%s.dbname", key)),
		User:   d.Config().GetKeyString(fmt.Sprintf("%s.user", key)),
		Pass:   d.Config().GetKeyString(fmt.Sprintf("%s.pass", key)),
		Host:   d.Config().GetKeyString(fmt.Sprintf("%s.host", key)),
		Port:   d.Config().GetKeyString(fmt.Sprintf("%s.port", key)),
	}

}

func (d *Dbs) Errors() []string {
	return d.errors
}

func (d *Dbs) AddError(e string) {
	d.errors = append(d.errors, e)
}
