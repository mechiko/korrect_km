package trueclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (t *trueClient) CisesListPost(target interface{}, cises []string) (string, error) {
	var u = t.urlGIS
	u.Path = `/api/v3/true-api/cises/short/list`
	body, err := json.Marshal(cises)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("%s %w", modError, err)
	}
	accept := "*/*"
	contentType := "application/json; charset=UTF-8"
	req.Header.Add("Accept", accept)
	req.Header.Add("Accept-Charset", "utf-8")
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.tokenGis))

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s %w", modError, err)
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	// t.Logger().Debugf("json_post:[%s]", buf)
	if resp.StatusCode != 200 {
		// 404 код не найден возвращает список всех кодов со статусом
		if resp.StatusCode != 404 {
			json.NewDecoder(bytes.NewBuffer(buf)).Decode(target)
			return string(buf), fmt.Errorf("ошибка запроса %d", resp.StatusCode)
		}
	}
	// потоковый Unmarshal
	return string(buf), json.NewDecoder(bytes.NewBuffer(buf)).Decode(target)
}
