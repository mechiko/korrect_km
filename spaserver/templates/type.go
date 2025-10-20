package templates

import (
	"html/template"
	"io"
	"io/fs"
	"korrectkm/domain"
	"path/filepath"

	"github.com/alexedwards/scs/v2"
	"github.com/mechiko/utility"
)

const modError = "http:templates"

const defaultTemplate = "index"

// если каталог ../spaserver/templates существует, то прописываем его в переменную
// для поиска шаблонов динамической обработки для отладки
// вычисляется абсолютный путь относительно каталога запуска в cmd под отладкой, потому ..
var pathTemplates = "../spaserver/templates"

func rootPathTemplates() (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = ""
		}
	}()
	out, err := filepath.Abs(pathTemplates)
	if err != nil {
		return ""
	}
	if utility.PathOrFileExists(out) {
		return out
	}
	return ""
}

type IApp interface {
	domain.Apper
	SessionManager() *scs.SessionManager
	SetTitlePage(string)
}

type ITemplateUI interface {
	LoadTemplates() (err error)
	Render(w io.Writer, page domain.Model, name string, data interface{}) error
	RenderDebug(w io.Writer, page domain.Model, name string, data interface{}) error
}

type Templates struct {
	IApp
	debug                    bool
	pages                    map[domain.Model]*template.Template
	fs                       fs.FS
	rootPathTemplateGinDebug string
	semaphore                Semaphore
}

var _ ITemplateUI = &Templates{}

// panic if error
func New(app IApp, debug bool) *Templates {
	t := &Templates{
		IApp:                     app,
		pages:                    nil,
		debug:                    debug,
		rootPathTemplateGinDebug: rootPathTemplates(),
		semaphore:                NewSemaphore(1),
	}
	if err := t.LoadTemplates(); err != nil {
		t.Logger().Errorf("%s %w", modError, err)
		panic(err.Error())
	}
	return t
}
