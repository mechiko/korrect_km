package templates

import (
	"fmt"
	"html/template"
	"io"
	"korrectkm/reductor"
	"path"
	"strings"
)

// debug
//   - true шаблоны грузятся каждый раз из файловой системы это для отладки
//   - false шаблоны парсятся при загрузке однажды
//
// все пути и включения отображаются из embeded структуры файлов, по ним строится t.pages[page]
// состоящая из дерева шаблонов для каждой страницы (независимых)
func (t *Templates) LoadTemplates() (err error) {
	t.pages = make(map[reductor.ModelType]*template.Template)
	t.fs = root
	embededPages, err := root.ReadDir(".")
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	for _, page := range embededPages {
		// t.Logger().Debugf("page %d %s %v", i, page.Name(), page.IsDir())
		if page.IsDir() {
			name := reductor.ModelTypeFromString(page.Name())
			if err := t.parsePage(name); err != nil {
				return fmt.Errorf("%s %w", modError, err)
			}
		}
	}
	return nil
}

func (t *Templates) parsePage(page reductor.ModelType) (err error) {
	// создаем новый шаблон страницы
	// при кэшировании мап не переписывается
	if _, ok := t.pages[page]; ok {
		return fmt.Errorf("%s такой шаблон вида %s уже обработан", modError, page)
	}
	t.pages[page] = template.New(page.String()).Funcs(functions)
	embededHtmls, err := root.ReadDir(page.String())
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	for _, html := range embededHtmls {
		// t.Logger().Debugf("htmls %d %s %v", i, html.Name(), html.IsDir())
		if !html.IsDir() {
			if err := t.parsePageHtml(page, html.Name(), t.pages[page]); err != nil {
				return fmt.Errorf("%s %w", modError, err)
			}
		}
	}
	return nil
}

func (t *Templates) parsePageHtml(page reductor.ModelType, html string, templ *template.Template) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()
	name, _ := strings.CutSuffix(path.Base(html), path.Ext(html))
	path := path.Join(page.String(), html)
	if file, err := t.fs.Open(path); err != nil {
		return fmt.Errorf("%s %w", modError, err)
	} else {
		if txt, err := io.ReadAll(file); err != nil {
			return fmt.Errorf("%s %w", modError, err)
		} else {
			templ.New(name).Funcs(functions).Parse(string(txt))
		}
	}
	return nil
}
