package stats

import (
	"bytes"
	"firstwails/utility"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (t *page) Route(e *echo.Echo) error {
	e.GET("/stats/title", t.Title)
	// e.GET("/stats/ready", t.Ready)
	e.GET("/stats/modal", t.modal)
	e.GET("/stats/file", t.file)
	e.GET("/stats/search", t.search)
	e.GET("/stats/reset", t.reset)
	e.GET("/stats/progress", t.progress)
	e.GET("/stats/excel/:status/:statusex", t.excel)
	e.GET("/stats/order/:status/:statusex", t.order)
	e.GET("/stats/upload", t.upload)
	e.POST("/stats/chunk", t.chunk)
	e.POST("/stats/orderload", t.orderload)
	e.POST("/stats/atkload", t.atkload)
	e.POST("/stats/utilload", t.utilload)
	e.GET("/stats/selectfile", t.selectFile)
	return nil
}

func (t *page) Title(c echo.Context) error {
	c.String(200, "STATS")
	return nil
}

// func (t *page) Ready(c echo.Context) error {
// 	if t.ReloadActivePage() {
// 		t.SetReloadActivePage(false)
// 		var buf bytes.Buffer
// 		model := t.Reductor().Model()
// 		if err := t.Render(&buf, "index", &model, c); err != nil {
// 			t.Logger().Errorf("%s %s", modError, err.Error())
// 			c.NoContent(204)
// 			return nil
// 		}
// 		c.HTML(200, buf.String())
// 		return nil
// 	}
// 	hx := htmxutil.HxRequestHeaderFromRequest(c.Request())
// 	if !hx.HxRequest {
// 		// готовность DOM когда не надо обновлять страницу и прилетел пинг от нее
// 		t.DomReadyHttp()
// 	}
// 	// без контента свап не производится
// 	c.NoContent(204)
// 	return nil
// }

func (t *page) modal(c echo.Context) error {
	var buf bytes.Buffer
	model := t.Reductor().Model()
	if err := t.Render(&buf, "modal", &model, c); err != nil {
		t.Logger().Errorf("%s %s", modError, err.Error())
		c.NoContent(204)
		return nil
	}
	c.HTML(200, buf.String())
	return nil
}

func (t *page) file(c echo.Context) error {
	file := utility.DialogOpenFile()
	t.Logger().Debugf("stats file selected: %s", file)
	if utility.PathOrFileExists(file) {
		var buf bytes.Buffer
		// обновляем модель Stats
		model := t.Reductor().Model()
		model.Stats.State = 1
		model.Stats.File = file
		t.UpdateModel(model, "stats.file")
		if err := t.Render(&buf, "page", &model, c); err != nil {
			t.Logger().Errorf("%s %s", modError, err.Error())
			c.NoContent(204)
			return nil
		}
		c.HTML(200, buf.String())
	}
	// возвращаем имя файла для отображения на форме
	// c.String(200, file)
	c.NoContent(204)
	return nil
}

func (t *page) search(c echo.Context) error {
	var buf bytes.Buffer
	// обрабатываем поиск
	model := t.Reductor().Model()
	model.Stats.CisOut = nil
	model.Stats.Errors = nil
	model.Stats.State = 1
	t.UpdateModel(model, "stats.search")
	// model = usecase.New(t).TrueClientSearch(model)
	t.Search(model)
	if err := t.Render(&buf, "page", &model, c); err != nil {
		t.Logger().Errorf("%s %s", modError, err.Error())
		c.NoContent(204)
		return nil
	}
	c.HTML(200, buf.String())
	return nil
}

func (t *page) reset(c echo.Context) error {
	var buf bytes.Buffer
	// обновляем модель Stats
	model := t.Reductor().Model()
	model.Stats.File = ""
	model.Stats.CisIn = nil
	model.Stats.CisOut = nil
	model.Stats.Errors = nil
	model.Stats.State = 0
	model.Stats.CisStatus = nil
	t.UpdateModel(model, "stats.reset")
	if err := t.Render(&buf, "page", &model, c); err != nil {
		t.Logger().Errorf("%s %s", modError, err.Error())
		c.NoContent(204)
		return nil
	}
	c.HTML(200, buf.String())
	return nil
}

func (t *page) progress(c echo.Context) error {
	var buf bytes.Buffer
	// обновляем модель Stats
	model := t.Reductor().Model().Stats
	if model.State == 1 {
		if err := t.Render(&buf, "progress", &model, c); err != nil {
			t.Logger().Errorf("%s %s", modError, err.Error())
			c.NoContent(204)
			return nil
		}
		c.String(200, buf.String())
	}
	c.String(200, "")
	return nil
}

func (t *page) excel(c echo.Context) error {
	status := c.Param("status")
	statusEx := c.Param("statusex")
	file := utility.DialogSaveFile()
	file = fmt.Sprintf("%s_%s", file, status)
	allcis := t.Reductor().Model().Stats.CisOut
	arrStatus := make([]string, 0)
	for _, cis := range allcis {
		if (cis.Status == status) && (cis.StatusEx == statusEx) {
			arrStatus = append(arrStatus, cis.Cis)
		}
	}
	size := t.Reductor().Model().Stats.ExcelChunkSize
	if size <= 0 {
		size = 30000
	}
	if err := t.ToExcel(arrStatus, file, size); err != nil {
		return c.String(200, err.Error())
	}
	t.Logger().Debugf("stats file selected: %s", file)
	// возвращаем имя файла для отображения на форме
	out := fmt.Sprintf("записано %d кодов маркировки", len(arrStatus))
	c.String(200, out)
	// c.NoContent(204)
	return nil
}

func (t *page) chunk(c echo.Context) error {
	chunkSize := c.FormValue("chunk")
	model := t.Reductor().Model()
	if size, err := strconv.Atoi(chunkSize); err != nil {
		model.Stats.ExcelChunkSize = 30000
	} else {
		model.Stats.ExcelChunkSize = size
	}
	t.UpdateModel(model, "stats.chunk")
	var buf bytes.Buffer
	if err := t.Render(&buf, "chunk", &model, c); err != nil {
		t.Logger().Errorf("%s %s", modError, err.Error())
		c.NoContent(204)
		return nil
	}
	c.HTML(200, buf.String())
	return nil
}

// создать заказ фиктивный для удобства обработки
func (t *page) order(c echo.Context) error {
	status := c.Param("status")
	statusEx := c.Param("statusex")
	allcis := t.Reductor().Model().Stats.CisOut
	arrStatus := make([]string, 0)
	for _, cis := range allcis {
		if (cis.Status == status) && (cis.StatusEx == statusEx) {
			arrStatus = append(arrStatus, cis.Cis)
		}
	}
	if len(arrStatus) < 1 {
		out := "с данным статусом кодов нет"
		return c.String(200, out)
	}
	cis := arrStatus[0]
	dbFile := t.Repo().Dbs().Lite().File()
	if id, err := t.Repo().DbZnak().AttachLite(dbFile, status, statusEx, cis); err != nil {
		out := "ошибка " + err.Error()
		return c.String(200, out)
	} else {
		out := fmt.Sprintf("создан заказ в АлкоХелп 3 №%d", id)
		return c.String(200, out)
	}
}

func (t *page) selectFile(c echo.Context) error {
	model := t.Reductor().Model()
	model.Stats.File = utility.DialogOpenFile()
	t.UpdateModel(model, "stats.selectFile")
	c.String(200, model.Stats.File)
	return nil
}
