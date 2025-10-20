package spaserver

import (
	"bytes"
	"fmt"
	"io"
	"korrectkm/domain"

	"github.com/labstack/echo/v4"
)

// пример отобрадения сообщения в обработке сервера
// view index это каталог в шаблонах spaserver\templates\index
// alert это файл шаблона spaserver\templates\index\alert.html
// if err := c.Render(http.StatusOK, "index", map[string]interface{}{
// 	"template": "alert",
// 	"data":     struct{ Error string }{Error: flushMsg},
// }); err != nil {
// 	s.ServerError(c, err)
// }

// name это имя view
// data это экземпляр map с именем шаблона для view и моделью данных содержащих два поля
// data это экземпляр map с именем шаблона для view и моделью данных содержащих два поля
// template имя файла в шаблонах вида без расширения и данные для шаблона data interface{}
//
//	map[string]interface{}{
//	            "template":   "Home Page",
//	            "data": domain.Model{},
//	        })
func (s *Server) Render(w io.Writer, name string, data interface{}, c echo.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()

	return s.renderToWriter(w, name, data)
}

// аналогичный рендеринг в строку из шаблона для использования в htmx
func (s *Server) RenderString(name string, data interface{}) (str string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic %v", r)
		}
	}()

	var buf bytes.Buffer
	if err := s.renderToWriter(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// renderToWriter handles the core rendering logic
func (s *Server) renderToWriter(w io.Writer, name string, data interface{}) error {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("data must be a map[string]interface{}")
	}
	if dataMap == nil {
		return fmt.Errorf("data map cannot be nil")
	}
	templateValue, exists := dataMap["template"]
	if !exists {
		return fmt.Errorf("data map must contain 'template' key")
	}
	subName, ok := templateValue.(string)
	if !ok {
		return fmt.Errorf("template value must be a string")
	}
	model := dataMap["data"]
	nameType, err := domain.ModelFromString(name)
	if err != nil {
		return fmt.Errorf("view name wrong string must be Application, TrueClient, StatusBar, NoPage, Header, Footer, Index, Home, Setup")
	}
	view, ok := s.views[nameType]
	if ok {
		switch nameType {
		case domain.Home:
			s.SetTitlePage(view.Title())
			// case reductor.Setup:
			// 	s.SetTitlePage(view.Title())
		}
	} else {
		s.Logger().Errorf("нет такого вида %s", name)
	}
	if s.debug {
		if err := s.templates.RenderDebug(w, nameType, subName, model); err != nil {
			s.Logger().Errorf("render debug %s", err.Error())
			return err
		}
	} else {
		if err := s.templates.Render(w, nameType, subName, model); err != nil {
			s.Logger().Errorf("render %s", err.Error())
			return err
		}
	}
	return nil
}
