package setup

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (t *page) AlertMessage(c echo.Context, model interface{}) {
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("alert", model)); err != nil {
		t.ServerError(c, err)
	}
}
