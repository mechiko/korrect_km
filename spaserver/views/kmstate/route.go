package kmstate

import (
	"errors"
	"fmt"
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
	t.Echo().GET(base+"/excel/:status/:statusex", t.excel)
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
	if err := t.ModelUpdate(data); err != nil {
		return t.ServerError(c, err)
	}
	go t.Search()
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("process_cis", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) excel(c echo.Context) error {
	status := c.Param("status")
	statusEx := c.Param("statusex")
	file, err := utility.DialogSaveFile(utility.Excel, "", ".")
	if err != nil {
		return t.ServerError(c, err)
	}
	file = fmt.Sprintf("%s_%s", file, status)
	data, err := t.PageModel()
	if err != nil {
		return t.ServerError(c, err)
	}
	allcis := data.CisOut
	arrStatus := make([]string, 0)
	for _, cis := range allcis {
		if (cis.Status == status) && (cis.StatusEx == statusEx) {
			arrStatus = append(arrStatus, cis.Cis)
		}
	}
	size := data.ExcelChunkSize
	if size <= 0 {
		size = 30000
	}
	fNames, err := t.ToExcel(arrStatus, file, size)
	if err != nil {
		return t.ServerError(c, err)
	}
	msg := fmt.Sprintf("записано %d КМ в %d файл(а,ов)", len(arrStatus), len(fNames))
	t.SetFlush(msg, "info")
	// возвращаем имя файла для отображения на форме
	// out := fmt.Sprintf("записано %d кодов маркировки %d файлов", len(arrStatus), len(fNames))
	// c.String(200, out)
	c.NoContent(204)
	return nil
}
