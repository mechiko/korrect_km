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
	Menu        Model = "menu"
	Home        Model = "home"
	Footer      Model = "footer"
	Root        Model = "root"
	KMState     Model = "kmstate"
	InnFias     Model = "innfias"
	Money       Model = "money"
)

func IsValidModel(s string) bool {
	switch Model(s) {
	case Application, TrueClient, StatusBar, NoPage, Header, Footer, Root, Home, KMState, Menu, InnFias, Money:
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
	case string(Menu):
		return Menu, nil
	case string(Footer):
		return Footer, nil
	case string(Home):
		return Home, nil
	case string(Root):
		return Root, nil
	case string(KMState):
		return KMState, nil
	case string(InnFias):
		return InnFias, nil
	case string(Money):
		return Money, nil
	}
	return "", fmt.Errorf("%s ошибочная модель domain.Model", s)
}

func (s Model) String() string {
	return string(s)
}
