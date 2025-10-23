package stats

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// загрузка КМ из заказа
func (t *page) orderload(c echo.Context) error {
	var bufRender bytes.Buffer
	orderid := c.FormValue("orderid")
	model := t.Reductor().Model()
	orderID, err := strconv.ParseInt(orderid, 10, 64)
	if err != nil {
		model.Stats.Errors = append(model.Stats.Errors, err.Error())
		t.UpdateModel(model, "stats.orderload.error")
		if err := t.Render(&bufRender, "page", &model, c); err != nil {
			t.Logger().Errorf("%s %s", modError, err.Error())
			c.NoContent(204)
			return nil
		}
		c.HTML(200, bufRender.String())
	}
	model.Stats.Errors = nil
	model.Stats.CisIn = nil
	model.Stats.File = ""
	model.Stats.CisStatus = make(map[string]map[string]int)
	model.Stats.State = 0
	if kms, err := t.Repo().DbZnak().OrderLoad(orderID); err != nil {
		model.Stats.Errors = append(model.Stats.Errors, err.Error())
		t.UpdateModel(model, "stats.orderload.error")
		if err := t.Render(&bufRender, "page", &model, c); err != nil {
			t.Logger().Errorf("%s %s", modError, err.Error())
			c.NoContent(204)
			return nil
		}
		c.HTML(200, bufRender.String())
	} else {
		for i := range kms {
			kms[i] = shrinkCisParse(kms[i])
		}
		model.Stats.CisIn = append(model.Stats.CisIn, kms...)
		// model.Stats.CisIn = utility.UniqueSliceElements(cisIn)
		t.UpdateModel(model, "stats.orderload")
		if err := t.Render(&bufRender, "page", &model, c); err != nil {
			t.Logger().Errorf("%s %s", modError, err.Error())
			c.NoContent(204)
			return nil
		}
		c.HTML(200, bufRender.String())
	}
	return nil
}

func shrinkCisParse(code string) string {
	index := strings.IndexByte(code, '\x1D')
	if index > 0 {
		return code[:index]
	}
	return code
}
