package views

import (
	"korrectkm/reductor"

	"github.com/labstack/echo/v4"
)

type IView interface {
	PageData() interface{}
	Routes() error
	Index(c echo.Context) error
	// имя подшаблона вида по умолчанию
	DefaultTemplate() string
	// имя подшаблона вида текущий
	CurrentTemplate() string
	// строковое значение reductor.ModelType
	Name() string
	// заголовок страницы
	Title() string
	InitData() interface{}
	ModelType() reductor.ModelType
	Svg() string
	Desc() string
}
