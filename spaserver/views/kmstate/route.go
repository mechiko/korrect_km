package kmstate

import (
	"korrectkm/reductor"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mechiko/utility"
)

func (t *page) Routes() error {
	// Serve static and media files under /static/ and /uploads/ path.
	base := "/" + t.modelType.String()
	t.Echo().GET(base, t.Index)
	t.Echo().GET(base+"/selectfile", t.selectFile)
	return nil
}

func (t *page) Index(c echo.Context) error {
	data, err := t.PageData()
	if err != nil {
		return t.ServerError(c, err)
	}
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("index", data)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) selectFile(c echo.Context) error {
	model, err := t.PageModel()
	if err != nil {
		return t.ServerError(c, err)
	}
	model.File, err = utility.DialogOpenFile([]utility.FileType{utility.Csv}, "", t.Pwd())
	if err != nil {
		return t.ServerError(c, err)
	}
	err = reductor.SetModel(model, false)
	if err != nil {
		return t.ServerError(c, err)
	}
	return c.String(200, model.File)
}
