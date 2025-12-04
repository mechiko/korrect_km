package menu

import (
	"korrectkm/domain"
	"korrectkm/spaserver/views"
	"strings"

	"github.com/labstack/echo/v4"
)

// const modError = "header"

type IApp interface {
	domain.Apper
	Echo() *echo.Echo
	ServerError(c echo.Context, err error) error
	Views() map[domain.Model]views.IView
	ActivePage() domain.Model
	SetActivePage(p domain.Model)
	Reload()
	Menu() []domain.Model
}

type page struct {
	IApp
	modelType       domain.Model
	defaultTemplate string
	currentTemplate string
	title           string
}

func New(app IApp) *page {
	t := &page{
		IApp:            app,
		modelType:       domain.Menu,
		defaultTemplate: "content",
		currentTemplate: "content",
	}
	return t
}

func (p *page) DefaultTemplate() string {
	return p.defaultTemplate
}

func (p *page) CurrentTemplate() string {
	return p.currentTemplate
}

func (p *page) Name() string {
	return strings.ToLower(p.modelType.String())
}

func (p *page) ModelType() domain.Model {
	return p.modelType
}
func (p *page) Title() string {
	return p.title
}

// иконка для меню
func (p *page) Svg() string {
	return ""
}

// описание вида для меню
func (p *page) Desc() string {
	return ""
}
