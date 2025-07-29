package header

import (
	"korrectkm/reductor"
)

type MenuModel struct {
	Title string
	Items MenuItemSlice
}

type MenuItem struct {
	Name   string
	Title  string
	Active bool
	Desc   string
	Svg    string
}

type MenuItemSlice []*MenuItem

func (t *page) InitData() interface{} {
	model := MenuModel{
		Title: "Меню",
		Items: make(MenuItemSlice, 0),
	}

	for _, m := range t.Menu() {
		page, ok := t.Views()[m]
		if !ok {
			t.Logger().Errorf("menu %s not found", m.String())
			continue
		}
		menuItem := &MenuItem{
			Name:   page.Name(),
			Title:  page.Title(),
			Active: t.ActivePage() == page.ModelType(),
			Desc:   page.Desc(),
			Svg:    page.Svg(),
		}
		model.Items = append(model.Items, menuItem)
	}
	reductor.Instance().SetModel(t.modelType, model)
	return model
}

func (t *page) PageData() interface{} {
	return reductor.Instance().Model(t.modelType)
}

// с преобразованием
func (t *page) PageModel() MenuModel {
	if mdl, ok := reductor.Instance().Model(t.modelType).(MenuModel); ok {
		return mdl
	}
	return MenuModel{}
}

// сброс модели редуктора для страницы
func (t *page) ResetData() {
}
