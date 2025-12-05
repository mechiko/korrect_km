package root

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (t *page) Routes() error {
	// Serve static and media files under /static/ and /uploads/ path.
	// t.Echo().GET("/", t.Index)
	// t.Echo().GET("/index.html", t.Index)
	return nil
}

func (t *page) Index(c echo.Context) error {
	data, err := t.PageData()
	if err != nil {
		return t.ServerError(c, err)
	}
	if err := c.Render(http.StatusOK, t.Name(), map[string]interface{}{
		"template": "index",
		"data":     data,
	}); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) Reset(c echo.Context) error {
	return nil
}
