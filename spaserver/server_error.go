package spaserver

import (
	"github.com/labstack/echo/v4"
)

// выводил лог ошибки и c.NoContent(204)
// в поток событий error sse
func (s *Server) ServerError(c echo.Context, err error) error {
	s.Logger().Errorf("%s server error %s", c.Request().RequestURI, err.Error())
	c.NoContent(204)
	s.SetFlush(err.Error(), "error")
	return err
}
