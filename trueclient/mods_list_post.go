package trueclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// target - domain.ModInfo
// param - domain.ModsListPostParam
func (t *trueClient) ModsListPost(target interface{}, param ModsListPostParam) error {
	var u = t.urlGIS
	u.Path = `/api/v3/true-api/mods/info`
	body, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	accept := "*/*"
	contentType := "application/json; charset=UTF-8"
	req.Header.Add("Accept", accept)
	req.Header.Add("Accept-Charset", "utf-8")
	req.Header.Add("Content-Type", contentType)
	// req.Header.Add("clientToken", t.tokenGis)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.tokenGis))
	// req.Header.Add("X-Signature", signBody)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("status %d %s", resp.StatusCode, buf)
	}
	t.Logger().Debugf("json:[%s]", buf)
	// потоковый Unmarshal
	return json.NewDecoder(bytes.NewBuffer(buf)).Decode(target)
}
