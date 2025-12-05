package header

import (
	"fmt"
	"korrectkm/domain"
)

type MenuModel struct {
	IApp
	Title string
	Items MenuItemSlice
	model domain.Model
}

type MenuItem struct {
	Name   string
	Title  string
	Active bool
	Desc   string
	Svg    string
}

type MenuItemSlice []*MenuItem

var _ domain.Modeler = (*MenuModel)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func NewModel(app IApp) (*MenuModel, error) {
	model := &MenuModel{
		IApp:  app,
		model: domain.Header,
		Title: "Меню",
		Items: make(MenuItemSlice, 0),
	}
	if err := model.ReadState(app); err != nil {
		return nil, fmt.Errorf("model MenuModel read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (m *MenuModel) SyncToStore(app domain.Apper) (err error) {
	return err
}

// читаем состояние приложения
func (m *MenuModel) ReadState(app domain.Apper) (err error) {
	return nil
}

func (m *MenuModel) Copy() (interface{}, error) {
	// shallow copy that`s why fields is simple
	dst := *m
	return &dst, nil
}

func (a *MenuModel) Model() domain.Model {
	return a.model
}

func (a *MenuModel) Save(_ domain.Apper) (err error) {
	return nil
}

func (a *MenuModel) SyncActive(active string) error {
	for _, m := range a.Items {
		m.Active = false
		if m.Name == active {
			m.Active = true
		}
	}
	return nil
}

func (a *MenuModel) ActiveTitle() string {
	for _, m := range a.Items {
		if m.Active {
			return m.Desc
		}
	}
	return ""
}
