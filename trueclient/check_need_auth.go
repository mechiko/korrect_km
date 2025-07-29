package trueclient

import "time"

// возвращает true если авторизация необходима
// без пинга СУЗ
func (t *trueClient) CheckNeedAuth() bool {
	if t.tokenGis == "" {
		return true
	}
	if t.tokenSuz == "" {
		return true
	}
	if t.authTime.IsZero() {
		return true
	}
	// если авторизация была больше 10 часов назад
	return time.Since(t.authTime) > 10*time.Hour
}

// возвращает true если авторизация необходима
// проверяется пингом СУЗ
func (t *trueClient) CheckNeedAuthPing() bool {
	if t.tokenGis == "" {
		t.pingSUZ = nil
		return true
	}
	if t.tokenSuz == "" {
		t.pingSUZ = nil
		return true
	}
	if png, err := t.PingSuz(); err != nil {
		// не важно что за ошибка, просто нужна авторизация
		t.Logger().Errorf("ошибка пинг при проверке необходимости авторизации ЧЗ %s", err.Error())
		return true
	} else {
		t.pingSUZ = png
	}
	return false
}
