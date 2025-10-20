package trueclient

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/domain/models/modeltrueclient"
	"net"
	"net/http"
)

// если не используется конфиг.дб то всегда совершает авторизацию и пинг
func NewFromModel(a domain.Apper, model *modeltrueclient.TrueClientModel) (s *trueClient, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", err)
		}
	}()

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			// Timeout: 5 * time.Second,
		}).Dial,
		// TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		// Timeout:   time.Second * 120,
		Transport: netTransport,
	}
	s = &trueClient{
		Apper:      a,
		httpClient: netClient,
		layout:     model.LayoutUTC,
		// logger:   zaplog.TrueSugar,
		urlSUZ:   model.StandSUZ,
		urlGIS:   model.StandGIS,
		tokenGis: model.TokenGIS,
		tokenSuz: model.TokenSUZ,
		hash:     model.HashKey,
		deviceID: model.DeviceID,
		omsID:    model.OmsID,
		authTime: model.AuthTime,
		errors:   make([]string, 0),
	}
	if (s.deviceID) == "" {
		// panic(fmt.Sprintf("%s %s", modError, "нужна настройка конфигурации"))
		return s, fmt.Errorf("%s %s", modError, "устройство должно быть заполнено")
	}
	if (s.omsID) == "" {
		// panic(fmt.Sprintf("%s %s", modError, "нужна настройка конфигурации"))
		return s, fmt.Errorf("%s %s", modError, "oms_id должно быть заполнено")
	}
	// проверяем необходимость авторизации пингом СУЗ
	if s.CheckNeedAuthPing() {
		// если авторизация необходима то проверяем использование конфиг.дб
		if model.IsConfigDB {
			// если нужна авторизация при использовании базы конфиг то ее надо делать в алкохелпе
			return s, fmt.Errorf("%s %s", modError, "необходимо получить токены авторизации в АлкоХелп 3")
		}
		if (s.hash) == "" {
			return s, fmt.Errorf("%s %s", modError, "ключ должен быть заполнен")
		}

		if err := s.AuthGisSuz(); err != nil {
			return s, fmt.Errorf("%s %s", modError, err.Error())
		}
		// сохраняем конфиг в объекте и редукторе
		s.Save(model)
		// сохраняем конфиг после авторизации
		model.SyncToStore(a)
	}
	return s, nil
}
