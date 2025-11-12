package home

import (
	"bytes"
	"image/png"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mechiko/barcode"
	"github.com/mechiko/barcode/datamatrix"
)

func (t *page) Routes() error {
	// Serve static and media files under /static/ and /uploads/ path.
	t.Echo().GET("/home", t.Index)
	t.Echo().GET("/datamatrix", t.Datamatrix)
	return nil
}

func (t *page) Index(c echo.Context) error {
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("index", t.PageModel())); err != nil {
		return t.ServerError(c, err)
	}
	t.SetFlush("папа у васи", "info")
	return nil
}

func (t *page) Datamatrix(c echo.Context) error {
	buffer := new(bytes.Buffer)
	code := "\xe80104630277410873215!,asF,l1k\"LH\x1d91EE11\x1d92PK6ejb9KiEm4jqt2G7tesaQ4bbukQfZumYfUrNxf9kE="
	dataMatrix, _ := datamatrix.Encode(code)
	qrCode, _ := barcode.Scale(dataMatrix, 100, 100)
	_ = png.Encode(buffer, qrCode)
	c.Blob(200, "image/png", buffer.Bytes())
	return nil
}

func (t *page) Reset(c echo.Context) error {
	return nil
}
