package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"korrectkm/domain"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// debug
//   - true шаблоны грузятся каждый раз из файловой системы это для отладки
//   - false шаблоны парсятся при загрузке однажды
//
// все пути и включения отображаются из embeded структуры файлов, по ним строится t.pages[page]
// состоящая из дерева шаблонов для каждой страницы (независимых)
func (t *Templates) DynLoadTemplates() (out map[domain.Model]*template.Template, err error) {
	out = make(map[domain.Model]*template.Template)
	fs := os.DirFS(t.rootPathTemplateGinDebug)
	embededPages, err := os.ReadDir(t.rootPathTemplateGinDebug)
	if err != nil {
		return out, fmt.Errorf("%s %w", modError, err)
	}
	for _, page := range embededPages {
		if page.IsDir() {
			name, err := domain.ModelFromString(page.Name())
			if err != nil {
				return out, fmt.Errorf("%s %w", modError, err)
			}
			if err := t.parsePageDyn(fs, name, out); err != nil {
				return out, fmt.Errorf("%s %w", modError, err)
			}
		}
	}
	return out, nil
}

func (t *Templates) parsePageDyn(fs fs.FS, page domain.Model, pages map[domain.Model]*template.Template) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()
	// создаем новый шаблон страницы
	pages[page] = template.New(page.String()).Funcs(functions)
	pg := page.String()
	path := filepath.Join(t.rootPathTemplateGinDebug, pg)
	embededHtmls, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	for _, html := range embededHtmls {
		if !html.IsDir() {
			if err := t.parsePageHtmlDyn(fs, page, html.Name(), pages[page]); err != nil {
				return fmt.Errorf("%s %w", modError, err)
			}
		}
	}
	return nil
}

func (t *Templates) parsePageHtmlDyn(fs fs.FS, page domain.Model, html string, templ *template.Template) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()

	name, _ := strings.CutSuffix(path.Base(html), path.Ext(html))
	pt := page.String()
	path := path.Join(pt, html)

	file, err := fs.Open(path)
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	defer file.Close()

	txt, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}

	if _, err := templ.New(name).Funcs(functions).Parse(string(txt)); err != nil {
		return fmt.Errorf("%s template parse error: %w", modError, err)
	}

	return nil
}
