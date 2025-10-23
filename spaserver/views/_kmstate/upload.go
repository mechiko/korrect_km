package stats

import (
	"bytes"
	"firstwails/utility"
	"fmt"
	"io"
	"os"

	"github.com/labstack/echo/v4"
)

// загрузка кодов КМ из файла формы
func (t *page) upload(c echo.Context) (errOut error) {
	var bufFile bytes.Buffer
	var bufRender bytes.Buffer

	model := t.Reductor().Model()
	model.Stats.Errors = nil
	model.Stats.CisIn = nil
	// model.Stats.File = ""
	model.Stats.CisStatus = make(map[string]map[string]int)
	model.Stats.State = 0
	if model.Stats.File == "" {
		errOut = fmt.Errorf("выберите файл")
	} else if !utility.PathOrFileExists(model.Stats.File) {
		errOut = fmt.Errorf("файл ненайден %s", model.Finder.File)
	} else if file, err := os.Open(model.Stats.File); err != nil {
		errOut = fmt.Errorf("ошибка файла %s", err.Error())
	} else {
		defer file.Close()
		if _, err = io.Copy(&bufFile, file); err != nil {
			errOut = fmt.Errorf("ошибка файла %s", err.Error())
		}
	}
	if errOut != nil {
		model.Stats.Errors = append(model.Stats.Errors, errOut.Error())
		t.UpdateModel(model, "stats.upload.error")
		if err := t.Render(&bufRender, "page", &model, c); err != nil {
			t.Logger().Errorf("%s %s", modError, err.Error())
			c.NoContent(204)
			return nil
		}
		c.HTML(200, bufRender.String())
	} else {
		cisIn := model.Stats.CisIn
		cisIn = append(cisIn, utility.ReadTextStringArrayReader(&bufFile)...)
		model.Stats.CisIn = utility.UniqueSliceElements(cisIn)
		t.UpdateModel(model, "stats.upload")
		if err := t.Render(&bufRender, "page", &model, c); err != nil {
			t.Logger().Errorf("%s %s", modError, err.Error())
			c.NoContent(204)
			return nil
		}
		c.HTML(200, bufRender.String())
	}
	return nil
}
