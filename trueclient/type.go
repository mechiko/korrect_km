package trueclient

import (
	"korrectkm/config"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// флаг запрещающий создание объекта изначально 0
var reentranceFlag int64

type ILogCfg interface {
	Config() config.IConfig
	Logger() *zap.SugaredLogger
}

const modError = "trueclient"

type trueClient struct {
	ILogCfg
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
	pingSUZ    *PingSuzInfo
}
