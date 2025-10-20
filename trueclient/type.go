package trueclient

import (
	"korrectkm/domain"
	"korrectkm/domain/models/modeltrueclient"
	"net/http"
	"net/url"
	"time"
)

// флаг запрещающий создание объекта изначально 0
var reentranceFlag int64

const modError = "trueclient"

type trueClient struct {
	domain.Apper
	urlSUZ url.URL
	urlGIS url.URL
	layout string
	// logger     *zap.SugaredLogger
	tokenGis   string // токен авторизации для урла
	tokenSuz   string
	hash       string // кэп
	deviceID   string
	omsID      string
	httpClient *http.Client
	authTime   time.Time
	errors     []string
	pingSUZ    *modeltrueclient.PingSuzInfo
}
