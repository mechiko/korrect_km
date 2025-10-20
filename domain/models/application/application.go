package application

import (
	"fmt"
	"korrectkm/config"
	"korrectkm/domain"
	"korrectkm/repo"
	"time"

	"github.com/mechiko/dbscan"
	"github.com/mechiko/utility"
)

type Application struct {
	model        domain.Model
	Title        string
	Export       string
	Browser      utility.Browser
	BrowserList  []string
	Output       string
	Debug        bool
	Host         string
	Port         string
	DbLiteDesc   string
	DbConfigDesc string
	DbZnakDesc   string
	DbA3Desc     string
	License      string
	FsrarID      string
	startTime    time.Time
	endTime      time.Time
	period       string
}

var _ domain.Modeler = (*Application)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func New(app domain.Apper) (*Application, error) {
	rp, err := repo.GetRepository()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	model := &Application{
		model:       domain.Application,
		Title:       "Application Title",
		BrowserList: []string{string(utility.Default), string(utility.Chrome), string(utility.Firefox), string(utility.Yandex), string(utility.Edge)},
	}
	if err := model.ReadState(app, rp); err != nil {
		return nil, fmt.Errorf("model application read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (m *Application) SyncToStore(app domain.Apper) (err error) {
	if err := app.SetOptions("export", m.Export); err != nil {
		return fmt.Errorf("application sync to store: set export failed: %w", err)
	}
	if err := app.SetOptions("browser", string(m.Browser)); err != nil {
		return fmt.Errorf("application sync to store: set browser failed: %w", err)
	}
	return nil
}

// читаем состояние приложения
func (m *Application) ReadState(app domain.Apper, rp *repo.Repository) (err error) {
	m.Export = app.Options().Export
	m.Browser = utility.Browser(app.Options().Browser)
	m.Output = app.Options().Output
	m.Host = app.Options().Hostname
	m.Port = app.Options().HostPort
	m.Debug = config.Mode == "development"
	for _, v := range rp.ListDbs() {
		if rp.Is(v) {
			info := rp.Info(v)
			if info == nil {
				// defensive: repo should guarantee non-nil, skip just in case
				continue
			}
			switch v {
			case dbscan.A3:
				switch info.Driver {
				case "sqlite":
					m.DbA3Desc = fmt.Sprintf("[%s] %s", info.Driver, info.File)
				case "mssql":
					m.DbA3Desc = fmt.Sprintf("[%s] %s", info.Driver, info.Name)
				default:
					m.DbA3Desc = fmt.Sprintf("[%s] %s", info.Driver, info.File)
				}
			case dbscan.Config:
				switch info.Driver {
				case "sqlite":
					m.DbConfigDesc = fmt.Sprintf("[%s] %s", info.Driver, info.File)
				case "mssql":
					m.DbConfigDesc = fmt.Sprintf("[%s] %s", info.Driver, info.Name)
				default:
					m.DbConfigDesc = fmt.Sprintf("[%s] %s", info.Driver, info.File)
				}
			case dbscan.TrueZnak:
				switch info.Driver {
				case "sqlite":
					m.DbZnakDesc = fmt.Sprintf("[%s] %s", info.Driver, info.File)
				case "mssql":
					m.DbZnakDesc = fmt.Sprintf("[%s] %s", info.Driver, info.Name)
				default:
					m.DbZnakDesc = fmt.Sprintf("[%s] %s", info.Driver, info.File)
				}
			case dbscan.Other:
				switch info.Driver {
				case "sqlite":
					m.DbLiteDesc = fmt.Sprintf("[%s] %s", info.Driver, info.File)
				case "mssql":
					m.DbLiteDesc = fmt.Sprintf("[%s] %s", info.Driver, info.Name)
				default:
					m.DbLiteDesc = fmt.Sprintf("[%s] %s", info.Driver, info.File)
				}
			default:
				return fmt.Errorf("application readstate type error %v", v)
			}
		}
	}
	m.License = app.Options().Application.License
	if cfgDb, err := rp.LockConfig(); err == nil {
		defer rp.UnlockConfig(cfgDb)
		if fsrarId, err := cfgDb.Key("fsrar_id"); err == nil {
			m.FsrarID = fsrarId
		}
	}
	return nil
}

func (a *Application) Copy() (interface{}, error) {
	// shallow copy that`s why fields is simple
	dst := *a
	return &dst, nil
}

func (a *Application) Model() domain.Model {
	return a.model
}

func (m *Application) Save(app domain.Apper) (err error) {
	if err := app.SaveOptions(); err != nil {
		return fmt.Errorf("application: save options failed: %w", err)
	}
	return nil
}
