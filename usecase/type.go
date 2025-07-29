package usecase

import (
	"korrectkm/config"
	"korrectkm/repo"

	"go.uber.org/zap"
)

const modError = "usecase"

type IApp interface {
	Logger() *zap.SugaredLogger
	Configuration() *config.Configuration
	Repo() *repo.Repository
}

type usecase struct {
	IApp
}

func New(app IApp) *usecase {
	return &usecase{
		IApp: app,
	}
}
