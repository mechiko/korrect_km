package trueclient

import (
	"korrectkm/domain/models/modeltrueclient"
	"time"
)

func (t *trueClient) TokenGIS() string {
	return t.tokenGis
}

func (t *trueClient) TokenSUZ() string {
	return t.tokenSuz
}

func (t *trueClient) AuthTime() time.Time {
	return t.authTime
}

func (t *trueClient) Errors() []string {
	return t.errors
}

func (t *trueClient) AddError(err error) {
	t.errors = append(t.errors, err.Error())
}

func (t *trueClient) PingSuzInfo() *modeltrueclient.PingSuzInfo {
	return t.pingSUZ
}
