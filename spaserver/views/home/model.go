package home

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/reductor"
)

type HomeModel struct {
	model   domain.Model
	Title   string
	CodeFNS string
}

var _ domain.Modeler = (*HomeModel)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func NewModel(app domain.Apper) (*HomeModel, error) {
	model := &HomeModel{
		model: domain.Home,
		Title: "HOME",
	}
	if err := model.ReadState(app); err != nil {
		return nil, fmt.Errorf("model %v read state %w", model.model, err)
	}
	return model, nil
}

// инициализируем модель вида
func (t *page) InitData(_ domain.Apper) (interface{}, error) {
	model := HomeModel{
		Title:   "HOME",
		CodeFNS: `0104630277410873215!,asF,l1k"LH91EE1192PK6ejb9KiEm4jqt2G7tesaQ4bbukQfZumYfUrNxf9kE=`,
	}
	err := reductor.SetModel(&model, false)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (m *HomeModel) SyncToStore(app domain.Apper) (err error) {
	return err
}

// читаем состояние приложения
func (m *HomeModel) ReadState(app domain.Apper) (err error) {
	return nil
}

func (m *HomeModel) Copy() (interface{}, error) {
	// shallow copy that`s why fields is simple
	dst := *m
	return &dst, nil
}

func (a *HomeModel) Model() domain.Model {
	return a.model
}

func (a *HomeModel) Save(_ domain.Apper) (err error) {
	return nil
}
