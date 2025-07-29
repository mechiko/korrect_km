package trueclient

import (
	"korrectkm/reductor"
)

// сохраняем в конфиге и в модели токены и время получения
func (t *trueClient) Save(model *TrueClientModel) {
	model.AuthTime = t.authTime
	model.TokenGIS = t.tokenGis
	model.TokenSUZ = t.tokenSuz
	mdl := *model
	reductor.Instance().SetModel(reductor.TrueClient, &mdl)
}
