package znaktool

import (
	"fmt"
	"time"
	"zupper/domain"
	"zupper/repo"
)

type ZnakTools struct {
	model       domain.Model
	Date        time.Time
	UtilNumber  int64
	OrderNumber int64
}

var _ domain.Modeler = (*ZnakTools)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func New(app domain.Apper, repo *repo.Repository) (*ZnakTools, error) {
	model := &ZnakTools{
		model: domain.ZnakTool,
		Date:  time.Now(),
	}
	if err := model.ReadState(app, repo); err != nil {
		return nil, fmt.Errorf("model ZnakArgegate read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения
func (m *ZnakTools) SyncToStore(app domain.Apper) (err error) {
	return err
}

// читаем состояние
func (m *ZnakTools) ReadState(app domain.Apper, repo *repo.Repository) (err error) {
	return nil
}

func (a *ZnakTools) Copy() (interface{}, error) {
	// shallow copy that`s why fields is simple
	dst := *a
	return &dst, nil
}

func (a *ZnakTools) Model() domain.Model {
	return a.model
}

func (a *ZnakTools) Save(_ domain.Apper) (err error) {
	return nil
}
