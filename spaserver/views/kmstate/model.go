package kmstate

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/reductor"
	"korrectkm/repo"

	"github.com/mechiko/dbscan"
)

type KmStateModel struct {
	model               domain.Model
	Title               string
	File                string
	CisIn               []string        // список CIS для запроса
	Chunks              int             // куски
	CisOut              domain.CisSlice // список CIS полученных
	CisStatus           map[string]map[string]int
	ExcelChunkSize      int      // размер куска для выгрузки в файл ексель
	IsConnectedTrueZnak bool     // есть подключение к ЧЗ
	IsTrueZnakA3        bool     // подключена БД ЧЗ А3
	AtkId               int      // номер ATK в ЧЗ А3
	OrderId             int      // номер заказа в ЧЗ А3
	UtilisationId       int      // номер отчета нанесения в ЧЗ А3
	Progress            int      // прогресс опроса
	IsProgress          bool     // true если идет процесс загрузки для отображения прогресса
	Errors              []string // массив ошибок
	MapCisStatusDict    map[string]string
	MapCisStatusExDict  map[string]string
}

var _ domain.Modeler = (*KmStateModel)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func NewModel(app domain.Apper) (*KmStateModel, error) {
	model := &KmStateModel{
		model:               domain.KMState,
		Title:               "Состояния КМ",
		IsConnectedTrueZnak: true, // подключаемся теперь при запуске приложения
		ExcelChunkSize:      30000,
	}
	if err := model.ReadState(app); err != nil {
		return nil, fmt.Errorf("model %v read state %w", model.model, err)
	}
	return model, nil
}

// инициализируем модель вида
func (t *page) InitData(app domain.Apper) (interface{}, error) {
	model, err := NewModel(app)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	err = reductor.SetModel(model, false)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения из модели редуктора
func (m *KmStateModel) SyncToStore(_ domain.Apper) (err error) {
	return err
}

// читаем состояние приложения
func (m *KmStateModel) ReadState(app domain.Apper) (err error) {
	rp, err := repo.GetRepository()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	m.IsTrueZnakA3 = rp.Is(dbscan.TrueZnak)
	return nil
}

func (m *KmStateModel) Copy() (interface{}, error) {
	// shallow copy that`s why fields is simple
	dst := *m
	return &dst, nil
}

func (a *KmStateModel) Model() domain.Model {
	return a.model
}

func (a *KmStateModel) Save(_ domain.Apper) (err error) {
	return nil
}
