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
	Application  Model = "application"
	TrueClient   Model = "trueclient"
	StatusBar    Model = "statusbar"
	ZnakAgregate Model = "znakagregate"
	ZnakTool     Model = "znaktool"
	NoPage       Model = "nopage"
	Header       Model = "header"
	Footer       Model = "footer"
	ProdTools    Model = "prodtools"
	Index        Model = "index"
)

func IsValidModel(s string) bool {
	switch Model(s) {
	case Application, TrueClient, StatusBar, ZnakAgregate, NoPage, Header, Footer, ProdTools, Index, ZnakTool:
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
	case string(ZnakAgregate):
		return ZnakAgregate, nil
	case string(ZnakTool):
		return ZnakTool, nil
	case string(NoPage):
		return NoPage, nil
	case string(Header):
		return Header, nil
	case string(Footer):
		return Footer, nil
	case string(ProdTools):
		return ProdTools, nil
	case string(Index):
		return Index, nil
	}
	return "", fmt.Errorf("%s ошибочная модель domain.Model", s)
}

func (s Model) String() string {
	return string(s)
}
