package spaserver

import (
	"fmt"
	"korrectkm/spaserver/views/footer"
	"korrectkm/spaserver/views/header"
	"korrectkm/spaserver/views/home"
	"korrectkm/spaserver/views/index"
	"net/http"

	"github.com/labstack/echo/v4"
)

// маршрутизация приложения
func (s *Server) Routes() http.Handler {
	s.loadViews()
	s.server.GET("/page", s.Page) // переход/загрузка на текущую страницу
	s.server.GET("/sse", s.Sse)
	return s.server
}

// вызывает рендеринг по имени страницы вида в s.activePage
func (s *Server) Page(c echo.Context) error {
	ap := s.activePage
	if _, ok := s.views[ap]; !ok {
		return s.ServerError(c, fmt.Errorf("not found active page %s", ap))
	}
	view := s.views[ap]
	return view.Index(c)
}

// загружаем все виды
func (s *Server) loadViews() {
	// view header
	view1 := header.New(s)
	s.views[view1.ModelType()] = view1
	view1.Routes()
	// view footer
	view2 := footer.New(s)
	s.views[view2.ModelType()] = view2
	view2.Routes()
	view2.InitData()
	// view home
	view3 := home.New(s)
	s.views[view3.ModelType()] = view3
	view3.Routes()
	view3.InitData()
	// view4 := setup.New(s)
	// s.views[view4.ModelType()] = view4
	// view4.Routes()
	// view4.InitData()
	// view index
	view5 := index.New(s)
	s.views[view5.ModelType()] = view5
	// header инициализируем последним нужны все виды сервера в списке
	view1.InitData()
}
