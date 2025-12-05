package views

import (
	"korrectkm/domain"

	"github.com/labstack/echo/v4"
)

type IView interface {
	// PageData() (interface{}, error)
	Routes() error
	Index(c echo.Context) error
	Reset(c echo.Context) error
	// имя подшаблона вида по умолчанию
	DefaultTemplate() string
	// имя подшаблона вида текущий
	CurrentTemplate() string
	// строковое значение reductor.ModelType
	Name() string
	// заголовок страницы
	Title() string
	InitData(domain.Apper) (interface{}, error)
	PageData() (interface{}, error)
	ModelType() domain.Model
	Svg() string
	Desc() string
}
