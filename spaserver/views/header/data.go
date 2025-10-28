package header

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/reductor"
)

func (t *page) InitData(app domain.Apper) (interface{}, error) {
	model, err := NewModel(app)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
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
	err = reductor.SetModel(model, false)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return model, nil
}

func (t *page) PageData() (interface{}, error) {
	return reductor.Model[*MenuModel](t.modelType)
}

// с преобразованием
func (t *page) PageModel() MenuModel {
	model, _ := reductor.Model[*MenuModel](t.modelType)
	return *model
}

// сброс модели редуктора для страницы
func (t *page) ResetData() {
}
