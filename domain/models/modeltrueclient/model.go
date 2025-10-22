package modeltrueclient

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/domain/mystore"
	"korrectkm/repo"
	"net/url"
	"time"

	"github.com/mechiko/dbscan"
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
	model := &TrueClientModel{
		model: domain.TrueClient,
	}
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
// когда считываем конфиг сбрасываем токены и время авторизации
func (m *TrueClientModel) ReadState(app domain.Apper) (err error) {
	m.TokenGIS = ""
	m.TokenSUZ = ""
	// time.IsZero()
	m.AuthTime = time.Time{}
	m.Test = app.Options().TrueClient.Test
	m.UseConfigDB = app.Options().TrueClient.UseConfigDB
	m.DeviceID = app.Options().TrueClient.DeviceID
	m.HashKey = app.Options().TrueClient.HashKey
	m.OmsID = app.Options().TrueClient.OmsID
	m.StandGIS = url.URL{
		Scheme: "https",
		Host:   app.Options().TrueClient.StandGIS,
	}
	if m.StandGIS.Host == "" {
		return fmt.Errorf("invalid or missing trueclient.standgis configuration")
	}
	m.StandSUZ = url.URL{
		Scheme: "https",
		Host:   app.Options().TrueClient.StandSUZ,
	}
	if m.StandSUZ.Host == "" {
		return fmt.Errorf("invalid or missing trueclient.standsuz configuration")
	}
	if m.Test {
		m.StandGIS = url.URL{
			Scheme: "https",
			Host:   app.Options().TrueClient.TestGIS,
		}
		m.StandSUZ = url.URL{
			Scheme: "https",
			Host:   app.Options().TrueClient.TestSUZ,
		}
	}
	rp, err := repo.GetRepository()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	m.IsConfigDB = rp.Is(dbscan.Config)
	if m.IsConfigDB && m.UseConfigDB {
		dbCfg, err := rp.LockConfig()
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		defer rp.UnlockConfig(dbCfg)

		m.OmsID, err = dbCfg.Key("oms_id")
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		m.DeviceID, err = dbCfg.Key("connection_id")
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		m.HashKey, err = dbCfg.Key("certificate_thumbprint")
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		m.TokenSUZ, err = dbCfg.Key("token_suz")
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		m.TokenGIS, err = dbCfg.Key("token_gis_mt")
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}

func (m *TrueClientModel) LoadStore(app domain.Apper) (err error) {
	// загружаем сертификаты пользователя
	m.MyStore, err = mystore.List(app.Logger())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
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
