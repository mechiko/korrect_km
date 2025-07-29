package setup

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// валидация без обновления модели редуктора
func (t *page) ValidateOmsID(c echo.Context) error {
	omsID := c.FormValue("omsid")
	model := map[string]interface{}{
		"template": "validate",
		"data":     struct{ Err string }{Err: ""},
	}
	if err := uuid.Validate(omsID); err != nil {
		model["data"] = struct{ Err string }{Err: err.Error()}
	}
	if out, err := t.RenderString("setup", model); err != nil {
		return t.ServerError(c, err)
	} else {
		h := t.Htmx().NewHandler(c.Response(), c.Request())
		if _, err := h.Write([]byte(out)); err != nil {
			return t.ServerError(c, err)
		}
	}
	return nil
}
