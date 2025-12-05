package znakdb

import (
	"fmt"
)

func (z *DbZnak) AtkCodes(id int64) (out []string, err error) {
	sess := z.dbSession
	codes := make([]map[string]interface{}, 0)
	res := sess.Collection("order_mark_atk_codes").Find("id_order_mark_atk", id)
	if err := res.All(&codes); err != nil {
		return nil, err
	}
	out = make([]string, len(codes))
	for i, code := range codes {
		cis, ok := code["code"].(string)
		if !ok {
			return nil, fmt.Errorf("%v not string %T", code["code"], code["code"])
		}
		out[i] = cis
	}
	return out, nil
}
