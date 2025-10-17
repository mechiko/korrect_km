package setup

import (
	"korrectkm/config"
	"korrectkm/reductor"
	"strings"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// const modError = "home"

type ILogCfg interface {
	Config() config.IConfig
	Logger() *zap.SugaredLogger
}

type IServer interface {
	Config() config.IConfig
	Logger() *zap.SugaredLogger
	Echo() *echo.Echo
	ServerError(c echo.Context, err error) error
	SetActivePage(reductor.ModelType)
	SetFlush(string, string)
	RenderString(name string, data interface{}) (str string, err error)
	Htmx() *htmx.HTMX
}

type page struct {
	IServer
	modelType       reductor.ModelType
	defaultTemplate string
	currentTemplate string
	title           string
}

func New(app IServer) *page {
	t := &page{
		IServer:         app,
		modelType:       reductor.Setup,
		defaultTemplate: "index",
		currentTemplate: "index",
		title:           "настройка соединения",
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

func (p *page) ModelType() reductor.ModelType {
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

// описание вида для меню
func (p *page) Desc() string {
	return "настройка чз"
}

func (p *page) ShowInMenu() bool {
	return true
}
