package setup

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/domain/models/modeltrueclient"
	"korrectkm/reductor"
	"korrectkm/repo"
	"korrectkm/trueclient"
	"time"

	"go.uber.org/zap"
)

type SetupModel struct {
	modeltrueclient.TrueClientModel
	PingSuz *modeltrueclient.PingSuzInfo
}

var _ domain.Modeler = (*SetupModel)(nil)

// берем данные из модели редуктора
func (m *SetupModel) Read(logger *zap.SugaredLogger) (err error) {
	return nil
}

func (m *SetupModel) ReadConfigDB() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	r, err := repo.GetRepository()
	if err != nil {
		return fmt.Errorf("setupModel repo error %w", err)
	}
	model, ok := reductor.Instance().Model(reductor.TrueClient).(trueclient.TrueClientModel)
	if !ok {
		return fmt.Errorf("setupModel read: модель не trueclient.TrueClientModel")
	}
	m.AuthTime = model.AuthTime
	m.DeviceID = model.DeviceID
	m.HashKey = model.HashKey
	m.IsConfigDB = model.IsConfigDB
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
		dbCfg, err := r.LockConfig()
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		defer r.UnlockConfig(dbCfg)
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

func (m *SetupModel) ClearConfigDB() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	model, ok := reductor.Instance().Model(domain.TrueClient).(trueclient.TrueClientModel)
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
	model, err := reductor.Instance().Model(domain.TrueClient)
	if err != nil {
		t.Logger().Errorf("setup: TrueClient error wrong %T", model)
		return SetupModel{}
	}
	return model
}

// с преобразованием
func (t *page) PageModel() SetupModel {
	model, err := reductor.Instance().Model(domain.TrueClient)
	if err != nil {
		t.Logger().Errorf("setup: TrueClient error wrong %T", model)
		return SetupModel{}
	}
	if mdl, ok := model.(SetupModel); ok {
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
