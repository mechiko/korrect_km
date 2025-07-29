package trueclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"korrectkm/trueclient/cmdsign"
	"net/http"
)

// target это куда возвращается json

func (t *trueClient) SearchGis(productGroups []string, gtins []string, target interface{}) error {
	u := t.urlGIS
	u.Path = `/api/v4/true-api/cises/search`
	t.Logger().Debugf("url:%s", u.String())
	body, err := t.filterSearchJson(productGroups, gtins)
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	t.Logger().Debugf("body: %s", body)
	signBody, err := cmdsign.New(t.hash).Sign(string(body))
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	accept := "application/json"
	req.Header.Add("Accept", accept)
	req.Header.Add("Content-Type", accept)
	// req.Header.Add("clientToken", t.tokenSuz)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.tokenGis))
	req.Header.Add("X-Signature", signBody)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s %w", modError, err)
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	t.Logger().Debugf("ping Body:%s", buf)
	// потоковый Unmarshal
	return json.NewDecoder(bytes.NewBuffer(buf)).Decode(target)
}

func (t *trueClient) filterSearchJson(productGroups []string, gtins []string) ([]byte, error) {
	filter := &filterSearch{
		ProductGroups: productGroups,
		Gtins:         gtins,
	}
	// 		ProductGroups: []string{
	// 			// "milk",
	// 			"beer",
	// 		},
	// 		Gtins: []string{"04810014009830"},
	// 		// OrderIds: []string{"9697ce05-66dc-4d5b-8b13-2f2fabfe545b"},
	// 	}
	// page := &paginationSearch{
	// 	PerPage:          1000,
	// 	LastEmissionDate: "2024-01-01T00:00:00",
	// 	Sgtin:            "0542500781002851TZN5V",
	// }
	sq := &searchQueryFilter{
		Filter: filter,
		// Pagination: page,
	}
	body, err := json.Marshal(sq)
	return body, err
}
