package stats

import (
	"bytes"
	"strconv"

	"github.com/labstack/echo/v4"
)

// загрузка КМ из заказа
func (t *page) utilload(c echo.Context) error {
	var bufRender bytes.Buffer
	orderid := c.FormValue("utilid")
	model := t.Reductor().Model()
	orderID, err := strconv.ParseInt(orderid, 10, 64)
	if err != nil {
		model.Stats.Errors = append(model.Stats.Errors, err.Error())
		t.UpdateModel(model, "stats.utilload.error")
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
	if kms, err := t.Repo().DbZnak().UtilisationLoad(orderID); err != nil {
		model.Stats.Errors = append(model.Stats.Errors, err.Error())
		t.UpdateModel(model, "stats.utilload.error")
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
		t.UpdateModel(model, "stats.utilload")
		if err := t.Render(&bufRender, "page", &model, c); err != nil {
			t.Logger().Errorf("%s %s", modError, err.Error())
			c.NoContent(204)
			return nil
		}
		c.HTML(200, bufRender.String())
	}
	return nil
}
