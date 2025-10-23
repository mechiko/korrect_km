package stats

import (
	"bytes"

	"github.com/labstack/echo/v4"
)

// загрузка КМ из заказа
func (t *page) atkload(c echo.Context) error {
	var bufRender bytes.Buffer
	atk := c.FormValue("atk")
	model := t.Reductor().Model()
	if atkID, err := t.Repo().DbZnak().FindATK(atk); err != nil {
		model.Stats.Errors = append(model.Stats.Errors, err.Error())
		t.UpdateModel(model, "stats.atkload.error")
		if err := t.Render(&bufRender, "page", &model, c); err != nil {
			t.Logger().Errorf("%s %s", modError, err.Error())
			c.NoContent(204)
			return nil
		}
		c.HTML(200, bufRender.String())
	} else {
		model.Stats.Errors = nil
		model.Stats.CisIn = nil
		model.Stats.File = ""
		model.Stats.CisStatus = make(map[string]map[string]int)
		model.Stats.State = 0
		if kms, err := t.Repo().DbZnak().AtkLoad(atkID); err != nil {
			model.Stats.Errors = append(model.Stats.Errors, err.Error())
			t.UpdateModel(model, "stats.atkload.error")
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
			t.UpdateModel(model, "stats.atkload")
			if err := t.Render(&bufRender, "page", &model, c); err != nil {
				t.Logger().Errorf("%s %s", modError, err.Error())
				c.NoContent(204)
				return nil
			}
			c.HTML(200, bufRender.String())
		}
	}
	return nil
}
