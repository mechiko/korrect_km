package repo

import (
	"errors"
	"fmt"
	"korrectkm/repo/a3"

	"github.com/mechiko/dbscan"
)

// if err is nil then must after Lock launch UnLock
// всегда или открывает базу и проверяет объект или возвращает ошибку
func (r *Repository) LockA3() (*a3.DbA3, error) {
	info := r.dbs.Info(dbscan.A3)
	if info == nil || !info.Exists {
		return nil, fmt.Errorf("%s lock info %v is nil or not exists", modError, dbscan.A3)
	}
	mu, ok := r.dbMutex[dbscan.A3]
	if ok {
		mu.mutex.Lock()
		// ensure we don't leak the lock on panic inside a3.New
		defer func() {
			if r := recover(); r != nil {
				mu.mutex.Unlock()
				panic(r)
			}
		}()
	} else {
		return nil, fmt.Errorf("%s lock not present mutex %v", modError, dbscan.A3)
	}
	db, err := a3.New(info)
	if err != nil {
		mu.mutex.Unlock()
		return nil, fmt.Errorf("%s lock not present mutex %v", modError, dbscan.A3)
	}
	return db, nil
}

func (r *Repository) UnlockA3(db *a3.DbA3) error {
	if db == nil {
		return fmt.Errorf("%s unlock db %v is nil", modError, dbscan.A3)
	}
	errClose := db.Close()
	mu, ok := r.dbMutex[db.InfoType()]
	if ok {
		mu.mutex.Unlock()
	} else {
		errUnlock := fmt.Errorf("%s unlock not present mutex %v", modError, db.InfoType())
		return errors.Join(errClose, errUnlock)
	}
	return errClose
}
