package index

import (
	"github.com/labstack/echo/v4"
)

func (t *page) Routes() error {
	// Serve static and media files under /static/ and /uploads/ path.
	// t.Echo().GET("/footer", t.Index)
	return nil
}

func (t *page) Index(c echo.Context) error {
	// if err := c.Render(http.StatusOK, t.Name(), map[string]interface{}{
	// 	"template": "content",
	// 	"data":     t.PageData(),
	// }); err != nil {
	// 	return t.ServerError(c, err)
	// }
	return nil
}

func (t *page) Reset(c echo.Context) error {
	return nil
}
