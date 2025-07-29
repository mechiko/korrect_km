package templates

import (
	"html/template"
	"io"
	"io/fs"
	"korrectkm/config"
	"korrectkm/reductor"

	"github.com/alexedwards/scs/v2"
	"go.uber.org/zap"
)

const modError = "http:templates"

const defaultTemplate = "index"

type IApp interface {
	Config() config.IConfig
	Logger() *zap.SugaredLogger
	SessionManager() *scs.SessionManager
	SetTitlePage(string)
	// DefaultPage() string
	// DefaultTemplate() string
	// ActivePage() string
	// ActiveTemplate() string
}

type ITemplateUI interface {
	LoadTemplates() (err error)
	Render(w io.Writer, page reductor.ModelType, name string, data interface{}) error
	RenderDebug(w io.Writer, page reductor.ModelType, name string, data interface{}) error
}

type Templates struct {
	IApp
	debug                    bool
	pages                    map[reductor.ModelType]*template.Template
	fs                       fs.FS
	rootPathTemplateGinDebug string
	semaphore                Semaphore
}

var _ ITemplateUI = &Templates{}

// panic if error
func New(app IApp, root string, debug bool) *Templates {
	t := &Templates{
		IApp:                     app,
		pages:                    nil,
		debug:                    debug,
		rootPathTemplateGinDebug: root,
		semaphore:                NewSemaphore(1),
	}
	if err := t.LoadTemplates(); err != nil {
		t.Logger().Errorf("%s %w", modError, err)
		panic(err.Error())
	}
	return t
}
