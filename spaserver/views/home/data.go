package home

import (
	"fmt"
	"korrectkm/reductor"
)

type HomeModel struct {
	Title   string
	CodeFNS string
}

// синхронизация с настройками если необходима
func (m *HomeModel) Sync(cfg ILogCfg) {
}

// чтение настроек если нужно
func (m *HomeModel) Read(cfg ILogCfg) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	return nil
}

func (t *page) PageData() interface{} {
	return reductor.Instance().Model(t.modelType)
}

// с преобразованием
func (t *page) PageModel() HomeModel {
	if mdl, ok := reductor.Instance().Model(t.modelType).(HomeModel); ok {
		return mdl
	}
	return HomeModel{}
}

// сброс модели редуктора для страницы
func (t *page) ResetData() {
}

func (t *page) ResetValidateData() {
}

// инициализируем модель вида
func (t *page) InitData() interface{} {
	model := HomeModel{
		Title:   "HOME",
		CodeFNS: `0104630277410873215!,asF,l1k"LH91EE1192PK6ejb9KiEm4jqt2G7tesaQ4bbukQfZumYfUrNxf9kE=`,
	}
	reductor.Instance().SetModel(t.modelType, model)
	return model
}
