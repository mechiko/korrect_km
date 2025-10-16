package configdb

import (
	"fmt"

	"github.com/upper/db/v4"
)

func (c *DbConfig) Key(k string) (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = fmt.Sprintf("panic %v", r)
		}
		c.Close()
	}()

	param := &Parameters{}
	if err := c.dbSession.Get(param, db.Cond{"name": k}); err != nil {
		return fmt.Sprintf("%s %v", modError, err)
	}
	return param.Value
}
