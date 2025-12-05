package spaserver

import (
	"fmt"
	"korrectkm/spaserver/views/footer"
	"korrectkm/spaserver/views/header"
	"korrectkm/spaserver/views/home"
	"korrectkm/spaserver/views/index"
	"korrectkm/spaserver/views/innfias"
	"korrectkm/spaserver/views/kmstate"
	"korrectkm/spaserver/views/menu"
	"korrectkm/spaserver/views/money"
	"net/http"

	"github.com/labstack/echo/v4"
)

// маршрутизация приложения
func (s *Server) Routes() http.Handler {
	s.loadViews()
	s.server.GET("/page", s.Page)        // переход/загрузка на текущую страницу
	s.server.GET("/page/reset", s.Reset) // переход/загрузка на текущую страницу
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

// вызывает рендеринг по имени страницы вида в s.activePage
func (s *Server) Reset(c echo.Context) error {
	ap := s.activePage
	if _, ok := s.views[ap]; !ok {
		return s.ServerError(c, fmt.Errorf("not found active page %s", ap))
	}
	view := s.views[ap]
	return view.Reset(c)
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
	view2.InitData(s)
	// view home
	view3 := home.New(s)
	s.views[view3.ModelType()] = view3
	view3.Routes()
	view3.InitData(s)
	// kmstate
	view4 := kmstate.New(s)
	s.views[view4.ModelType()] = view4
	view4.Routes()
	view4.InitData(s)
	// view index
	view5 := index.New(s)
	s.views[view5.ModelType()] = view5
	// menu
	view6 := menu.New(s)
	s.views[view6.ModelType()] = view6
	view6.Routes()
	// innfias
	view7 := innfias.New(s)
	s.views[view7.ModelType()] = view7
	view7.Routes()
	view7.InitData(s)
	// money
	view8 := money.New(s)
	s.views[view8.ModelType()] = view8
	view8.Routes()
	view8.InitData(s)
	// header инициализируем последним нужны все виды сервера в списке
	view1.InitData(s)
	view6.InitData(s)
}
