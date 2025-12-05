package home

import (
	_ "embed"
	"korrectkm/domain"
	"korrectkm/spaserver/views"
	"strings"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

// const modError = "home"

type IServer interface {
	domain.Apper
	Echo() *echo.Echo
	ServerError(c echo.Context, err error) error
	SetActivePage(domain.Model)
	SetFlush(string, string)
	RenderString(name string, data interface{}) (str string, err error)
	Htmx() *htmx.HTMX
}

type page struct {
	IServer
	modelType domain.Model
	// name            string
	defaultTemplate string
	currentTemplate string
	title           string
}

var _ views.IView = (*page)(nil)

func New(app IServer) *page {
	t := &page{
		IServer:         app,
		modelType:       domain.Home,
		defaultTemplate: "index",
		currentTemplate: "index",
		title:           "домашняя страница",
	}
	return t
}

// шаблон по умолчанию это на будущее
func (p *page) DefaultTemplate() string {
	return p.defaultTemplate
}

// текущий шаблон это на будущее
func (p *page) CurrentTemplate() string {
	return p.currentTemplate
}

// low caps name
func (p *page) Name() string {
	return strings.ToLower(p.modelType.String())
}

func (p *page) ModelType() domain.Model {
	return p.modelType
}

// формируем мап для рендера map[string]interface{}{template": .., "data"...}
func (p *page) RenderPageModel(tmpl string, model interface{}) map[string]interface{} {
	return map[string]interface{}{
		"template": tmpl,
		"data":     model,
	}
}

func (p *page) Title() string {
	return p.title
}

//go:embed home.svg
var svg string

// иконка для меню
func (p *page) Svg() string {
	return svg
}

// описание вида для меню
func (p *page) Desc() string {
	return "домой"
}
