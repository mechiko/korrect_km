package app

import (
	"context"
	"fmt"

	"korrectkm/config"
	"korrectkm/reductor"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ILogCfg interface {
	Config() config.IConfig
	Logger() *zap.SugaredLogger
}

type app struct {
	ctx       context.Context
	uuid      string // идентификатор для уникальности формы
	config    config.IConfig
	logger    *zap.SugaredLogger
	pwd       string
	startTime time.Time
	endTime   time.Time
	output    string
}

// var _ ILogCfg = &webapp{}
var _ ILogCfg = (*app)(nil)

const modError = "app"

// func NewWebApp(logger *zap.SugaredLogger, e *echo.Echo, sse *sse.Server, pwd string) *webapp {
func NewWebApp(cfg config.IConfig, logger *zap.SugaredLogger, pwd string) *app {
	sc := &app{}
	sc.pwd = pwd
	sc.logger = logger
	sc.config = cfg
	sc.uuid = uuid.New().String()
	logger.Info("start pages")
	sc.initDateMn()
	model := ApplicationModel{
		Title:   "Application Title",
		License: "",
		FsrarID: "",
	}
	model.Read(sc)
	reductor.Instance().SetModel(reductor.Application, model)
	return sc
}

func (a *app) initDateMn() {
	loc, _ := time.LoadLocation("Europe/Moscow")
	t := time.Now().In(loc)
	_, m, _ := t.Date()
	a.startTime = time.Date(t.Year(), m, 1, 1, 0, 0, 0, loc)
	a.endTime = a.startTime.AddDate(0, 1, -1)
}

func (a *app) NowDateString() string {
	n := time.Now()
	return fmt.Sprintf("%4d.%02d.%02d %02d:%02d:%02d", n.Local().Year(), n.Local().Month(), n.Local().Day(), n.Local().Hour(), n.Local().Minute(), n.Local().Second())
}

func (a *app) StartDateString() string {
	return fmt.Sprintf("%4d.%02d.%02d", a.startTime.Local().Year(), a.startTime.Local().Month(), a.startTime.Local().Day())
}

func (a *app) EndDateString() string {
	return fmt.Sprintf("%4d.%02d.%02d", a.endTime.Local().Year(), a.endTime.Local().Month(), a.endTime.Local().Day())
}

func (a *app) SetStartDate(d time.Time) {
	a.startTime = d
}

func (a *app) SetEndDate(d time.Time) {
	a.endTime = d
}

func (a *app) StartDate() time.Time {
	return a.startTime
}

func (a *app) EndDate() time.Time {
	return a.endTime
}

func (a *app) FsrarID() string {
	return a.Config().Configuration().Application.Fsrarid
}

func (a *app) SetFsrarID(id string) {
	a.Config().SetInConfig("application.fsrarid", id, true)
}

func (a *app) Pwd() string {
	return a.pwd
}

func (a *app) Output() string {
	return a.output
}

func (a *app) Config() config.IConfig {
	return a.config
}

func (a *app) Logger() *zap.SugaredLogger {
	return a.logger
}

func (a *app) Ctx() context.Context {
	return a.ctx
}
