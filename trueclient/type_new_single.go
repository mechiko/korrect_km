package trueclient

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/domain/models/modeltrueclient"
	"net"
	"net/http"
	"sync/atomic"
)

// если не используется конфиг.дб то всегда совершает авторизацию
// блокирует реентерабельность до выполнения ClearSingle()
func NewFromModelSingle(a domain.Apper, model *modeltrueclient.TrueClientModel) (s *trueClient, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w", err)
			if s != nil {
				s.errors = append(s.errors, err.Error())
			}
		}
	}()

	if !atomic.CompareAndSwapInt64(&reentranceFlag, 0, 1) {
		errStr := "установлен запрет на создание объекта"
		return nil, fmt.Errorf("%s", errStr)
	}
	defer func() {
		atomic.StoreInt64(&reentranceFlag, 0)
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
		urlSUZ:     model.StandSUZ,
		urlGIS:     model.StandGIS,
		tokenGis:   model.TokenGIS,
		tokenSuz:   model.TokenSUZ,
		hash:       model.HashKey,
		deviceID:   model.DeviceID,
		omsID:      model.OmsID,
		authTime:   model.AuthTime,
	}

	if (s.deviceID) == "" {
		return s, fmt.Errorf("%s %s", modError, "устройство должно быть заполнено")
	}
	if (s.omsID) == "" {
		return s, fmt.Errorf("%s %s", modError, "oms_id должно быть заполнено")
	}
	// проверяем необходимость авторизации пингом СУЗ
	if s.CheckNeedAuthPing() {
		// если авторизация необходима то проверяем использование конфиг.дб
		if model.UseConfigDB && model.IsConfigDB {
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
		// проверяем авторизацию пингом
		if _, err := s.PingSuz(); err != nil {
			// не важно что за ошибка, просто нужна авторизация
			return s, fmt.Errorf("%s %s", modError, err.Error())
		}
	}
	model.PingSuz = s.pingSUZ
	return s, nil
}
