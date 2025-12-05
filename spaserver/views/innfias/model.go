package innfias

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/reductor"
)

type InnFiasModel struct {
	model domain.Model
	Title string
}

var _ domain.Modeler = (*InnFiasModel)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func NewModel(app domain.Apper) (*InnFiasModel, error) {
	model := &InnFiasModel{
		model: domain.InnFias,
		Title: "Информация по ИНН",
	}
	if err := model.ReadState(app); err != nil {
		return nil, fmt.Errorf("model %v read state %w", model.model, err)
	}
	return model, nil
}

// инициализируем модель вида
func (t *page) InitData(app domain.Apper) (interface{}, error) {
	model, err := NewModel(app)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	err = reductor.SetModel(model, false)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (m *InnFiasModel) SyncToStore(_ domain.Apper) (err error) {
	return err
}

// читаем состояние приложения
func (m *InnFiasModel) ReadState(app domain.Apper) (err error) {
	return nil
}

func (m *InnFiasModel) Copy() (interface{}, error) {
	// shallow copy that`s why fields is simple
	dst := *m
	return &dst, nil
}

func (a *InnFiasModel) Model() domain.Model {
	return a.model
}

func (a *InnFiasModel) Save(_ domain.Apper) (err error) {
	return nil
}
