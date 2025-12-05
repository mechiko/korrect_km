package menu

import (
	"fmt"
	"korrectkm/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (t *page) Routes() error {
	t.Echo().GET("/menu", t.Index)
	t.Echo().GET("/menu/:page", t.pager)
	return nil
}

func (t *page) Index(c echo.Context) error {
	data, err := t.PageData()
	if err != nil {
		return t.ServerError(c, err)
	}
	if err := c.Render(http.StatusOK, t.Name(), map[string]interface{}{
		"template": "content",
		"data":     data,
	}); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) modal(c echo.Context) error {
	model, err := t.InitData(t)
	if err != nil {
		return t.ServerError(c, err)
	}
	if err := c.Render(http.StatusOK, t.Name(), map[string]interface{}{
		"template": "modal",
		"data":     model,
	}); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) pager(c echo.Context) error {
	page := c.Param("page")
	if page == "" {
		return t.ServerError(c, fmt.Errorf("pager page empty"))
	}
	tm, err := domain.ModelFromString(page)
	if err != nil {
		return t.ServerError(c, err)
	}
	model, err := t.PageModel()
	if err != nil {
		return t.ServerError(c, err)
	}
	t.SetActivePage(tm)
	model.SyncActive(page)
	err = t.ModelUpdate(model)
	if err != nil {
		return t.ServerError(c, err)
	}
	t.Reload()
	return nil
}

func (t *page) Reset(c echo.Context) error {
	return nil
}
