package znakagregate

import (
	"fmt"
	"zupper/domain"
	"zupper/repo"
)

const defaultItemPerGroup = 6

type ZnakAgregate struct {
	model domain.Model
	// GroupOrders          domain.OrderInfoSlice
	// PackageOrders        domain.OrderInfoSlice
	// SelectedGroupOrder   *domain.OrderInfo
	// SelectedPackageOrder *domain.OrderInfo
	ItemPerGroup int
	FileName     string
}

var _ domain.Modeler = (*ZnakAgregate)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func New(app domain.Apper, repo *repo.Repository) (*ZnakAgregate, error) {
	model := &ZnakAgregate{
		model:        domain.ZnakAgregate,
		ItemPerGroup: defaultItemPerGroup,
	}
	if err := model.ReadState(app, repo); err != nil {
		return nil, fmt.Errorf("model ZnakArgegate read state %w", err)
	}
	// if !reductor.Instance().IsExistModel(model.model) {
	// 	if err := reductor.Instance().SetModel(model.model, model); err != nil {
	// 		return nil, fmt.Errorf("model %s store to reductor error %v", model.model, err)
	// 	}
	// }
	return model, nil
}

// синхронизирует с приложением в сторону приложения
func (m *ZnakAgregate) SyncToStore(_ domain.Apper) (err error) {
	return err
}

// читаем состояние
func (m *ZnakAgregate) ReadState(_ domain.Apper, _ *repo.Repository) (err error) {
	return nil
}

func (a *ZnakAgregate) Copy() (interface{}, error) {
	// shallow copy that`s why fields is simple
	dst := *a
	return &dst, nil
}

func (a *ZnakAgregate) Model() domain.Model {
	return a.model
}

func (a *ZnakAgregate) Save(_ domain.Apper) (err error) {
	return nil
}
