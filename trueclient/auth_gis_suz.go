package trueclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"korrectkm/trueclient/cmdsign"
	"net/http"
	"net/url"
	"path"
	"time"
)

const authPath = `/api/v3/true-api/auth/key`
const signGIS = `/api/v3/true-api/auth/simpleSignIn`
const signSUZ = `/api/v3/true-api/auth/simpleSignIn`

func (t *trueClient) AuthGisSuz() (err error) {
	authJSON := struct {
		Uuid string `json:"uuid"`
		Data string `json:"data"`
		Inn  string `json:"inn"`
	}{}
	// две попытки подключения
	// почему то в режиме отладки почте всегда первая попытка неудачная ...
	// не знаю почему
	// attempt := 2
	// for {
	// 	err := t.getAuth(authPath, &authJSON)
	// 	if err == nil {
	// 		break
	// 	}
	// 	attempt--
	// 	if attempt == 0 {
	// 		return fmt.Errorf("%w", err)
	// 	}
	// }
	err = t.getAuth(authPath, &authJSON)
	if err != nil {
		// t.logger.Errorf("getAuth:%s", err.Error())
		return fmt.Errorf("getAuth %w", err)
	}
	t.Logger().Infof("authJSON:%+v", authJSON)
	authJSON.Data, err = cmdsign.New(t.hash).Sign(authJSON.Data)
	if err != nil {
		// t.logger.Errorf("cmdsign:%s", err.Error())
		return fmt.Errorf("cmdsign.New %w", err)
	}
	body, err := json.Marshal(authJSON)
	if err != nil {
		// t.logger.Errorf("json.Marshal(authJSON):%s", err.Error())
		return fmt.Errorf("%w", err)
	}
	tokenJSON := struct {
		Token string `json:"token"`
	}{}
	err = t.postSignGis(signGIS, body, &tokenJSON)
	if err != nil {
		// t.logger.Errorf("json.Marshal(authJSON):%s", err.Error())
		return fmt.Errorf("postSignGis %w", err)
	}
	t.Logger().Infof("postSignGis:%.6s…%.4s", tokenJSON.Token, tokenJSON.Token)
	t.tokenGis = tokenJSON.Token

	err = t.postSignSuz(signSUZ, body, &tokenJSON)
	if err != nil {
		// t.logger.Errorf("postSignSuz:%s", err.Error())
		return fmt.Errorf("postSignSuz %w", err)
	}
	t.Logger().Infof("tokenSuz:%.6s…%.4s", tokenJSON.Token, tokenJSON.Token)
	t.tokenSuz = tokenJSON.Token
	t.authTime = time.Now()
	return nil
}

// target - адрес структуры для разбора JSON
func (t *trueClient) getAuth(path string, target interface{}) (err error) {
	var u = url.URL{
		Scheme: t.urlGIS.Scheme,
		Host:   t.urlGIS.Host,
		Path:   path,
	}
	r, err := t.httpClient.Get(u.String())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer r.Body.Close()
	buf, _ := io.ReadAll(r.Body)
	if r.StatusCode != 200 {
		return fmt.Errorf("%s", buf)
	}
	// потоковый Unmarshal
	return json.NewDecoder(bytes.NewBuffer(buf)).Decode(target)
}

func (t *trueClient) postSignGis(path string, body []byte, target interface{}) error {
	var u = url.URL{
		Scheme: t.urlGIS.Scheme,
		Host:   t.urlGIS.Host,
		Path:   path,
	}
	contentType := "application/json"
	// data := []byte(body)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	req.Header.Add("Content-Type", contentType)
	// req.Header.Add("Authorization", "Bearer YOUR_ACCESS_TOKEN")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", buf)
	}

	return json.NewDecoder(bytes.NewBuffer(buf)).Decode(target)
}

func (t *trueClient) postSignSuz(pathStr string, body []byte, target interface{}) error {
	var u = url.URL{
		Scheme: t.urlGIS.Scheme,
		Host:   t.urlGIS.Host,
		Path:   path.Join(pathStr, t.deviceID),
	}
	contentType := "application/json"
	// data := []byte(body)
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	req.Header.Add("Content-Type", contentType)
	// req.Header.Add("Authorization", "Bearer YOUR_ACCESS_TOKEN")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", buf)
	}
	return json.NewDecoder(bytes.NewBuffer(buf)).Decode(target)
}
