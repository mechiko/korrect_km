package app

import (
	"korrectkm/utility"

	"github.com/skratchdot/open-golang/open"
)

func (a *app) OpenDir() {
	defer func() {
		if r := recover(); r != nil {
			a.Logger().Errorf("%s panic %v", modError, r)
		}
	}()

	if a.Config().Configuration().Output == "" {
		return
	}
	if err := open.Run(a.Config().Configuration().Output); err != nil {
		a.Logger().Errorf("Dir %s %s", a.Config().Configuration().Output, err.Error())
	}
}

func (a *app) Open(url string) {
	if url == "" {
		return
	}
	if err := utility.OpenURL(url); err != nil {
		a.Logger().Errorf("%w", err)
	}
}
