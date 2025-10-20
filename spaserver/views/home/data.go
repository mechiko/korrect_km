package home

import (
	"korrectkm/domain"
	"korrectkm/reductor"
)

func (t *page) PageData() interface{} {
	mdl, _ := reductor.Instance().Model(domain.Home)
	return mdl
}

// с преобразованием
func (t *page) PageModel() HomeModel {
	model, _ := reductor.Instance().Model(domain.Header)
	if mdl, ok := model.(HomeModel); ok {
		return mdl
	}
	return HomeModel{}
}

// сброс модели редуктора для страницы
func (t *page) ResetData() {
}

func (t *page) ResetValidateData() {
}
