package trueclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Host:   "markirovka.sandbox.crptech.ru"
func (t *trueClient) Balance(target interface{}, productId int64) error {
	var u = url.URL{
		Scheme: t.urlGIS.Scheme,
		Host:   t.urlGIS.Host,
		Path:   `/api/v3/true-api/elk/product-groups/balance/all`,
	}
	if productId != 0 {
		productGroupId := strconv.FormatInt(productId, 10)
		var v = make(url.Values)
		v.Set("productGroupId", productGroupId)
		u.Path = `/api/v3/true-api/elk/product-groups/balance`
		u.RawQuery = v.Encode()
	} else {
		// query = u.Path
	}
	t.Logger().Debugf("url:%s", u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	accept := "application/json"
	contentType := "application/json; charset=UTF-8"
	req.Header.Add("Accept", accept)
	req.Header.Add("Accept-Charset", "utf-8")
	req.Header.Add("Content-Type", contentType)
	// req.Header.Add("clientToken", t.tokenGis)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.tokenGis))
	// req.Header.Add("X-Signature", querySign)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	t.Logger().Debugf("balance Body:%s", buf)
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s %s", modError, buf)
	}
	// потоковый Unmarshal
	return json.NewDecoder(bytes.NewBuffer(buf)).Decode(target)
}
