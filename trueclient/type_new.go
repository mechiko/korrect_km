package trueclient

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/domain/models/modeltrueclient"
	"korrectkm/reductor"
	"net"
	"net/http"
	"sync/atomic"
)

// инициализируем структурой с полями
// проверка необходимиости авторизации и ее выполнение
// model указатель и изменяется авторизацией
func New(a domain.Apper) (s *trueClient, err error) {
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

	// флаг устанавливаем в методе создания только одного объекта
	if !atomic.CompareAndSwapInt64(&reentranceFlag, 0, 1) {
		return nil, fmt.Errorf("установлен запрет на создание объекта")
	}

	// при запуске программы модель должна быть инициализирована
	// здесь мы уже получаем ее существующую
	model, err := reductor.Model[*modeltrueclient.TrueClientModel](domain.TrueClient)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
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
	validateField := func(value, fieldName string) error {
		if value == "" {
			return fmt.Errorf("%s %s: %s", modError, "нужна настройка конфигурации", fieldName)
		}
		return nil
	}
	if err := validateField(s.deviceID, "deviceId"); err != nil {
		return nil, err
	}
	if err := validateField(s.deviceID, "omsId"); err != nil {
		return nil, err
	}
	if err := validateField(s.deviceID, "hash"); err != nil {
		return nil, err
	}
	// проверяем необходимость авторизации пингом СУЗ
	if s.CheckNeedAuthPing() {
		if model.IsConfigDB {
			// если нужна авторизация при использовании базы конфиг то ее надо делать в алкохелпе
			return nil, fmt.Errorf("%s %s", modError, "необходимо получить токены авторизации в АлкоХелп 3")
		}
		if err := s.AuthGisSuz(); err != nil {
			return nil, fmt.Errorf("%s %w", modError, err)
		}
		// сохраняем конфиг в объекте
		s.Save(model)
		// сохраняем конфиг после авторизации
		model.SyncToStore(s)
	}
	return s, nil
}
