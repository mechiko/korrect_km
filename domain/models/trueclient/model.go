package trueclient

import (
	"fmt"
	"korrectkm/domain"
	"net/url"
	"time"
)

type PingSuzInfo struct {
	OmsId      string `json:"omsId"`
	ApiVersion string `json:"apiVersion"`
	OmsVersion string `json:"omsVersion"`
}

func (p *PingSuzInfo) String() string {
	return fmt.Sprintf("OMS ID:%s\nAPI:%s\nOMS:%s\n", p.OmsId, p.ApiVersion, p.OmsVersion)
}

type TrueClientModel struct {
	model       domain.Model
	Title       string
	StandGIS    url.URL
	StandSUZ    url.URL
	TokenGIS    string
	TokenSUZ    string
	AuthTime    time.Time
	LayoutUTC   string
	HashKey     string
	DeviceID    string
	OmsID       string
	IsConfigDB  bool // есть ли база конфиг.дб выставляется при запуске
	UseConfigDB bool // если ли база конфиг.дб есть то копируем данные из нее для авторизации
	Errors      []string
	PingSuz     *PingSuzInfo
	Validates   map[string]string
	MyStore     map[string]string
	Test        bool
}

var _ domain.Modeler = (*TrueClientModel)(nil)

// создаем модель считываем ее состояние и возвращаем указатель
func New(app domain.Apper) (*TrueClientModel, error) {
	model := &TrueClientModel{}
	if err := model.ReadState(app); err != nil {
		return nil, fmt.Errorf("model ZnakArgegate read state %w", err)
	}
	return model, nil
}

// синхронизирует с приложением в сторону приложения
func (m *TrueClientModel) SyncToStore(app domain.Apper) (err error) {
	app.SetOptions("trueclient.deviceid", m.DeviceID)
	app.SetOptions("trueclient.hashkey", m.HashKey)
	app.SetOptions("trueclient.omsid", m.OmsID)
	app.SetOptions("trueclient.useconfigdb", m.UseConfigDB)
	return err
}

// читаем состояние
func (m *TrueClientModel) ReadState(_ domain.Apper) (err error) {
	return nil
}

func (a *TrueClientModel) Copy() (interface{}, error) {
	// shallow copy that`s why fields is simple
	dst := *a
	return &dst, nil
}

func (a *TrueClientModel) Model() domain.Model {
	return a.model
}

func (a *TrueClientModel) Save(_ domain.Apper) (err error) {
	return nil
}
