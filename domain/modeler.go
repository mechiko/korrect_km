package domain

import (
	"fmt"
	"strings"
)

type Modeler interface {
	Save(Apper) error
	Copy() (interface{}, error) // структура копирует себя и выдает ссылку на копию с массивами и другими данными
	Model() Model               // возвращает тип модели
}

type Model string

const (
	Application Model = "application"
	TrueClient  Model = "trueclient"
	StatusBar   Model = "statusbar"
	NoPage      Model = "nopage"
	Header      Model = "header"
	Home        Model = "home"
	Footer      Model = "footer"
	Index       Model = "index"
	Setup       Model = "setup"
)

func IsValidModel(s string) bool {
	switch Model(s) {
	case Application, TrueClient, StatusBar, NoPage, Header, Footer, Index, Home, Setup:
		return true
	default:
		return false
	}
}

// строка приводится в нижний регистр потом сравнивается
func ModelFromString(s string) (Model, error) {
	s = strings.ToLower(s)
	switch s {
	case string(Application):
		return Application, nil
	case string(TrueClient):
		return TrueClient, nil
	case string(StatusBar):
		return StatusBar, nil
	case string(NoPage):
		return NoPage, nil
	case string(Header):
		return Header, nil
	case string(Footer):
		return Footer, nil
	case string(Home):
		return Home, nil
	case string(Index):
		return Index, nil
	case string(Setup):
		return Setup, nil
	}
	return "", fmt.Errorf("%s ошибочная модель domain.Model", s)
}

func (s Model) String() string {
	return string(s)
}
