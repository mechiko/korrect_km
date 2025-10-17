package checkdbg

import (
	"fmt"

	"korrectkm/domain"
)

const modError = "pkg:checkdbg"

type Checks struct {
	domain.Apper
}

func NewChecks(app domain.Apper) *Checks {
	return &Checks{
		Apper: app,
	}
}

func (c *Checks) Run() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s Run panic %v", modError, r)
		}
	}()

	return nil
}
