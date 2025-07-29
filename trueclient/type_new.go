package trueclient

import (
	"fmt"
	"korrectkm/reductor"
	"net"
	"net/http"
	"sync/atomic"
)

// инициализируем структурой с полями
// проверка необходимиости авторизации и ее выполнение
// model указатель и изменяется авторизацией
func New(a ILogCfg) (s *trueClient) {
	s = &trueClient{
		errors: make([]string, 0),
	}
	defer func() {
		if r := recover(); r != nil {
			errStr := fmt.Sprintf("%s panic %v", modError, r)
			s.errors = append(s.errors, errStr)
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

	// флаг устанавливаем в методе создания только одного объекта
	if !atomic.CompareAndSwapInt64(&reentranceFlag, 0, 1) {
		s.errors = append(s.errors, "установлен запрет на создание объекта")
		return s
	}

	// при запуске программы модель должна быть инициализирована
	// здесь мы уже получаем ее существующую
	model, ok := reductor.Instance().Model(reductor.TrueClient).(TrueClientModel)
	if !ok {
		strErr := fmt.Sprintf("reductor model trueclient wrong type %T", reductor.Instance().Model(reductor.TrueClient))
		a.Logger().Errorf("reductor model trueclient wrong type %T", reductor.Instance().Model(reductor.TrueClient))
		panic(strErr)
	}
	s = &trueClient{
		ILogCfg:    a,
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
	validateField := func(value, fieldName string) {
		if value == "" {
			panic(fmt.Sprintf("%s %s: %s", modError, "нужна настройка конфигурации", fieldName))
		}
	}
	validateField(s.deviceID, "deviceId")
	validateField(s.omsID, "omsId")
	validateField(s.hash, "hash")
	// проверяем необходимость авторизации пингом СУЗ
	if s.CheckNeedAuthPing() {
		if model.IsConfigDB {
			// если нужна авторизация при использовании базы конфиг то ее надо делать в алкохелпе
			panic(fmt.Sprintf("%s %s", modError, "необходимо получить токены авторизации в АлкоХелп 3"))
		}
		if err := s.AuthGisSuz(); err != nil {
			panic(fmt.Sprintf("%s %s", modError, err.Error()))
		}
		// сохраняем конфиг в объекте
		s.Save(&model)
		// сохраняем конфиг после авторизации
		model.Sync(s)
	}
	return s
}
