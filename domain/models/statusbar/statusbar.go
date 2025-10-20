package statusbar

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/repo"
)

type StatusBar struct {
	model   domain.Model
	Utm     bool
	License bool
	Scan    bool
	FsrarID string
}

var _ domain.Modeler = (*StatusBar)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func New(app domain.Apper, repo *repo.Repository) (*StatusBar, error) {
	model := &StatusBar{
		model: domain.StatusBar,
	}
	if err := model.ReadState(app, repo); err != nil {
		return nil, fmt.Errorf("model statusbar read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (m *StatusBar) SyncToStore(app domain.Apper) (err error) {
	return err
}

// читаем состояние приложения
func (m *StatusBar) ReadState(app domain.Apper, repo *repo.Repository) (err error) {
	m.Utm = false
	m.License = true
	m.Scan = false
	m.FsrarID = app.Options().Application.Fsrarid
	return nil
}

func (a *StatusBar) Copy() (interface{}, error) {
	// shallow copy that`s why fields is simple
	dst := *a
	return &dst, nil
}

func (a *StatusBar) Model() domain.Model {
	return a.model
}

func (a *StatusBar) Save(_ domain.Apper) (err error) {
	return nil
}
