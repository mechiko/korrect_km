package checkdbg

import (
	"fmt"

	"korrectkm/config"

	"go.uber.org/zap"
)

const modError = "pkg:checkdbg"

type ILogCfg interface {
	Config() config.IConfig
	Logger() *zap.SugaredLogger
}

type Checks struct {
	ILogCfg
}

func NewChecks(app ILogCfg) *Checks {
	return &Checks{
		ILogCfg: app,
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
