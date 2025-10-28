package trueclient

import (
	"korrectkm/domain/models/modeltrueclient"
	"korrectkm/reductor"
)

// сохраняем в конфиге и в модели токены и время получения
func (t *trueClient) Save(model *modeltrueclient.TrueClientModel) {
	model.AuthTime = t.authTime
	model.TokenGIS = t.tokenGis
	model.TokenSUZ = t.tokenSuz
	if err := reductor.SetModel(model, false); err != nil {
		t.Logger().Errorf("trueclient save redutctor model error %w", err)
	}
}
