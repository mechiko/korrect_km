package reductor

import (
	"sync"

	"go.uber.org/zap"
)

type Reductor struct {
	mutex  sync.Mutex
	logger *zap.SugaredLogger
	models ModelList
}

var once sync.Once
var instance *Reductor

// func New(pageModel ModelType, model IModel, logger *zap.SugaredLogger) *Reductor {
func New(logger *zap.SugaredLogger) *Reductor {
	once.Do(func() {
		instance = &Reductor{
			logger: logger,
			models: make(ModelList),
		}
	})
	return instance
}

func Instance() *Reductor {
	return instance
}

func (rdc *Reductor) Logger() *zap.SugaredLogger {
	return rdc.logger
}
