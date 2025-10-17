package a3

import (
	"fmt"
)

// вызывается каждый раз при создании объекта
func (r *DbA3) Check() (err error) {
	if !r.dbInfo.Exists {
		return fmt.Errorf("dbConfig dbInfo.Exists false")
	}
	r.dbSession, err = r.dbInfo.Connect()
	if err != nil {
		return fmt.Errorf("%s check ошибка подключения к БД %w", modError, err)
	}
	return nil
}
