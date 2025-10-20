package trueclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"korrectkm/domain/models/modeltrueclient"
	"net/http"
	"net/url"
)

// возвращает true если пинг есть
func (t *trueClient) PingSuzSilent() bool {
	if t.tokenSuz == "" {
		return false
	}
	var v = make(url.Values)
	v.Set("omsId", t.omsID)
	var u = url.URL{
		Scheme:   t.urlSUZ.Scheme,
		Host:     t.urlSUZ.Host,
		Path:     `/api/v3/ping`,
		RawQuery: v.Encode(),
	}
	t.Logger().Debugf("url:%s", u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return false
	}
	accept := "application/json"
	req.Header.Add("Accept", accept)
	req.Header.Add("Content-Type", accept)
	req.Header.Add("clientToken", t.tokenSuz)
	// req.Header.Add("X-Signature", querySign)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	t.Logger().Debugf("ping Body:%s", buf)
	// потоковый Unmarshal
	pingJSON := modeltrueclient.PingSuzInfo{}
	if err := json.NewDecoder(bytes.NewBuffer(buf)).Decode(&pingJSON); err != nil {
		return false
	}
	t.pingSUZ = &pingJSON
	return true
}

// https://suz.sandbox.crptech.ru/api/v3/ping?omsId=32539e31-c671-4462-8443-3d92b038b0f9
func (t *trueClient) PingSuz() (info *modeltrueclient.PingSuzInfo, err error) {
	info = &modeltrueclient.PingSuzInfo{}
	var v = make(url.Values)
	v.Set("omsId", t.omsID)
	var u = url.URL{
		Scheme:   t.urlSUZ.Scheme,
		Host:     t.urlSUZ.Host,
		Path:     `/api/v3/ping`,
		RawQuery: v.Encode(),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return info, fmt.Errorf("%s %w", modError, err)
	}
	accept := "application/json"
	req.Header.Add("Accept", accept)
	req.Header.Add("Content-Type", accept)
	req.Header.Add("clientToken", t.tokenSuz)
	// req.Header.Add("X-Signature", querySign)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return info, fmt.Errorf("%s %w", modError, err)
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return info, fmt.Errorf("%s %s", modError, buf)
	}
	// потоковый Unmarshal
	err = json.NewDecoder(bytes.NewBuffer(buf)).Decode(info)
	t.pingSUZ = info
	return info, err
}
