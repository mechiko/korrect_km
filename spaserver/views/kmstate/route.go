package kmstate

import (
	"errors"
	"korrectkm/reductor"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mechiko/utility"
)

func (t *page) Routes() error {
	// Serve static and media files under /static/ and /uploads/ path.
	base := "/" + t.modelType.String()
	t.Echo().GET(base, t.Index)
	t.Echo().GET(base+"/selectfile", t.selectFile)
	t.Echo().GET(base+"/reset", t.Reset)
	t.Echo().GET(base+"/progress", t.progress)
	t.Echo().GET(base+"/search", t.search)
	return nil
}

func (t *page) Index(c echo.Context) error {
	data, err := t.PageData()
	if err != nil {
		return t.ServerError(c, err)
	}
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("index", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

// сброс к начальному состоянию
func (t *page) Reset(c echo.Context) error {
	data, err := t.InitData(t)
	if err != nil {
		return t.ServerError(c, err)
	}
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("index", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) selectFile(c echo.Context) error {
	model, err := t.PageModel()
	if err != nil {
		return t.ServerError(c, err)
	}
	model.File, err = utility.DialogOpenFile([]utility.FileType{utility.Csv}, "", t.Pwd())
	if err != nil {
		return t.ServerError(c, err)
	}
	err = model.loadCisFromFile()
	if err != nil {
		return t.ServerError(c, err)
	}
	err = reductor.SetModel(model, false)
	if err != nil {
		return t.ServerError(c, err)
	}
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("page", model)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) progress(c echo.Context) error {
	data, err := t.PageModel()
	if err != nil {
		return t.ServerError(c, err)
	}
	if data.IsProgress {
		if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("progress", data)); err != nil {
			return t.ServerError(c, err)
		}
	} else {
		if len(data.Errors) > 0 {
			err := strings.Join(data.Errors, "<BR>")
			data.Errors = make([]string, 0)
			t.ModelUpdate(data)
			return t.ServerError(c, errors.New(err))
		} else {
			if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("page_stop", data)); err != nil {
				return t.ServerError(c, err)
			}
		}
	}
	return nil
}

func (t *page) search(c echo.Context) error {
	data, err := t.PageModel()
	if err != nil {
		return t.ServerError(c, err)
	}
	if data.IsProgress {
		// уже запущена операция
		// return c.NoContent(204)
		return t.ServerError(c, errors.New("уже запущена операция, дождитесь ее окончания"))
	}
	// запускаем прогресс чтобы отобразить на странице он отображается когда больше 0
	data.IsProgress = true
	go t.Search()
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("process_cis", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}
