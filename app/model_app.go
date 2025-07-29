package app

import (
	"fmt"
)

type ApplicationModel struct {
	Title   string // заголовок главного окна будет в конфиге не хранится
	License string
	FsrarID string
}

// что храним в конфиге тут прописываем
func (m *ApplicationModel) Sync(cfg ILogCfg) {
	cfg.Config().SetInConfig("application.license", m.License, true)
	cfg.Config().SetInConfig("application.fsrarid", m.FsrarID, true)
}

// что читаем из конфига тут
func (m *ApplicationModel) Read(cfg ILogCfg) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	m.License = cfg.Config().GetKeyString("application.license")
	m.FsrarID = cfg.Config().GetKeyString("application.fsrarid")
	return nil
}
