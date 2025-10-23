package kmstate

import (
	"korrectkm/reductor"
)

func (t *page) PageData() (interface{}, error) {
	return reductor.Instance().Model(t.modelType)
}

// с преобразованием
func (t *page) PageModel() KmStateModel {
	model, _ := reductor.Instance().Model(t.modelType)
	if mdl, ok := model.(KmStateModel); ok {
		return mdl
	}
	return KmStateModel{}
}

// сброс модели редуктора для страницы
func (t *page) ResetData() {
}

func (t *page) ResetValidateData() {
}
