package header

import (
	"korrectkm/domain"
	"korrectkm/reductor"
)

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
	reductor.Instance().SetModel(&model, false)
	return model
}

func (t *page) PageData() interface{} {
	mdl, _ := reductor.Instance().Model(domain.Header)
	return mdl
}

// с преобразованием
func (t *page) PageModel() MenuModel {
	model, _ := reductor.Instance().Model(domain.Header)
	if mdl, ok := model.(MenuModel); ok {
		return mdl
	}
	return MenuModel{}
}

// сброс модели редуктора для страницы
func (t *page) ResetData() {
}
