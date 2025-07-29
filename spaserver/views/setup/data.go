package setup

import (
	"fmt"
	"korrectkm/reductor"
	"korrectkm/repo"
	"korrectkm/trueclient"
	"korrectkm/trueclient/mystore"
	"net/url"
	"time"

	"go.uber.org/zap"
)

type SetupModel struct {
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
	UseConfigDB bool
	IsConfigDB  bool
	Errors      []string
	PingSuz     *trueclient.PingSuzInfo
	Validates   map[string]string
	MyStore     map[string]string
	Test        bool
}

// синхронизация с настройками и с моделью редуктора TrueClient если необходима
func (m *SetupModel) Sync(cfg ILogCfg) {
	cfg.Config().SetInConfig("trueclient.test", m.Test, true)
	cfg.Config().SetInConfig("trueclient.deviceid", m.DeviceID, true)
	cfg.Config().SetInConfig("trueclient.hashkey", m.HashKey, true)
	cfg.Config().SetInConfig("trueclient.omsid", m.OmsID, true)
	cfg.Config().SetInConfig("trueclient.useconfigdb", m.UseConfigDB, true)
	tclModel, _ := m.ToTrueClient()
	reductor.Instance().SetModel(reductor.TrueClient, tclModel)

}

// берем данные из модели редуктора
func (m *SetupModel) Read(logger *zap.SugaredLogger) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	model, ok := reductor.Instance().Model(reductor.TrueClient).(trueclient.TrueClientModel)
	if !ok {
		return fmt.Errorf("setupModel read: модель не trueclient.TrueClientModel")
	}
	m.AuthTime = model.AuthTime
	m.DeviceID = model.DeviceID
	m.HashKey = model.HashKey
	m.IsConfigDB = model.IsConfigDB
	m.UseConfigDB = model.UseConfigDB
	m.LayoutUTC = model.LayoutUTC
	m.MyStore = model.MyStore
	m.OmsID = model.OmsID
	m.PingSuz = model.PingSuz
	m.StandGIS = model.StandGIS
	m.StandSUZ = model.StandSUZ
	m.Test = model.Test
	m.TokenGIS = model.TokenGIS
	m.TokenSUZ = model.TokenSUZ
	if m.MyStore, err = mystore.List(logger); err != nil {
		logger.Errorf("%s", err.Error())
	}
	return nil
}

func (m *SetupModel) ReadConfigDB(r *repo.Repository) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	model, ok := reductor.Instance().Model(reductor.TrueClient).(trueclient.TrueClientModel)
	if !ok {
		return fmt.Errorf("setupModel read: модель не trueclient.TrueClientModel")
	}
	m.AuthTime = model.AuthTime
	m.DeviceID = model.DeviceID
	m.HashKey = model.HashKey
	m.IsConfigDB = r.IsConfig()
	m.UseConfigDB = true
	m.LayoutUTC = model.LayoutUTC
	m.MyStore = model.MyStore
	m.OmsID = model.OmsID
	m.PingSuz = model.PingSuz
	m.StandGIS = model.StandGIS
	m.StandSUZ = model.StandSUZ
	m.Test = model.Test
	m.TokenGIS = model.TokenGIS
	m.TokenSUZ = model.TokenSUZ
	if m.IsConfigDB {
		m.OmsID = r.ConfigDB().Key("oms_id")
		m.DeviceID = r.ConfigDB().Key("connection_id")
		m.HashKey = r.ConfigDB().Key("certificate_thumbprint")
		m.TokenSUZ = r.ConfigDB().Key("token_suz")
		m.TokenGIS = r.ConfigDB().Key("token_gis_mt")
	}
	return nil
}

func (m *SetupModel) ClearConfigDB() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	model, ok := reductor.Instance().Model(reductor.TrueClient).(trueclient.TrueClientModel)
	if !ok {
		return fmt.Errorf("setupModel read: модель не trueclient.TrueClientModel")
	}
	m.AuthTime = time.Time{}
	m.DeviceID = model.DeviceID
	m.HashKey = model.HashKey
	m.IsConfigDB = model.IsConfigDB
	m.UseConfigDB = false
	m.LayoutUTC = model.LayoutUTC
	m.MyStore = model.MyStore
	m.OmsID = model.OmsID
	m.PingSuz = nil
	m.StandGIS = model.StandGIS
	m.StandSUZ = model.StandSUZ
	m.Test = model.Test
	m.TokenGIS = ""
	m.TokenSUZ = ""
	return nil
}

func (t *page) PageData() interface{} {
	if mdl, ok := reductor.Instance().Model(t.modelType).(SetupModel); ok {
		// обновляем при запросе данных с моделью труклиент
		mdl.Read(t.Logger())
		return mdl
	}
	return reductor.Instance().Model(t.modelType)
}

// с преобразованием
func (t *page) PageModel() SetupModel {
	if mdl, ok := reductor.Instance().Model(t.modelType).(SetupModel); ok {
		mdl.Read(t.Logger())
		return mdl
	}
	return SetupModel{}
}

// сброс модели редуктора для страницы
func (t *page) ResetData() {
}

func (t *page) ResetValidateData() {
}

// инициализируем модель вида
func (t *page) InitData() interface{} {
	model := SetupModel{
		Title: "Настройка соединения",
	}
	model.Read(t.Logger())
	reductor.Instance().SetModel(t.modelType, model)
	return model
}

// берем данные из модели редуктора
func (m *SetupModel) ToTrueClient() (model trueclient.TrueClientModel, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	model = trueclient.TrueClientModel{}
	model.AuthTime = m.AuthTime
	model.DeviceID = m.DeviceID
	model.HashKey = m.HashKey
	model.IsConfigDB = m.IsConfigDB
	model.UseConfigDB = m.UseConfigDB
	model.LayoutUTC = m.LayoutUTC
	model.MyStore = m.MyStore
	model.OmsID = m.OmsID
	model.PingSuz = m.PingSuz
	model.StandGIS = m.StandGIS
	model.StandSUZ = m.StandSUZ
	model.Test = m.Test
	model.TokenGIS = m.TokenGIS
	model.TokenSUZ = m.TokenSUZ
	return model, nil
}
