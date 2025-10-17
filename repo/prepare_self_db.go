package repo

import (
	"fmt"
	"korrectkm/repo/selfdb"

	"github.com/mechiko/dbscan"
)

// при инициализации приложения этот метод вызывается однажды и прописывает объект доступа
// к базе данных, далее проверяет версию БД возможна ошибка и нужно выходить из приложения
func (r *Repository) prepareSelf() (err error) {
	defer func() {
		if rr := recover(); rr != nil {
			err = fmt.Errorf("%s panic %v", modError, rr)
		}
	}()
	selfInfo := r.Info(dbscan.Other)
	if selfInfo == nil {
		//return fmt.Errorf("%s get info: nil", modError)
		return fmt.Errorf("%s lock info %v is nil or not exists", modError, dbscan.Other)
	}
	self, err := selfdb.New(selfInfo)
	if err != nil {
		return fmt.Errorf("%s self new error %w", modError, err)
	}
	defer self.Close()

	if !selfInfo.Exists {
		return fmt.Errorf("%s get self db error!", modError)
	}
	return nil
}
