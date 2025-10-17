package repo

import (
	"errors"
	"fmt"
	"korrectkm/repo/znakdb"

	"github.com/mechiko/dbscan"
)

// if err is nil then must after Lock launch UnLock
// всегда или открывает базу и проверяет объект или возвращает ошибку
func (r *Repository) LockZnak() (*znakdb.DbZnak, error) {
	info := r.dbs.Info(dbscan.TrueZnak)
	if info == nil || !info.Exists {
		return nil, fmt.Errorf("%s lock info %v is nil or not exists", modError, dbscan.TrueZnak)
	}
	mu, ok := r.dbMutex[dbscan.TrueZnak]
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
		return nil, fmt.Errorf("repo lock not present mutex %v", dbscan.TrueZnak)
	}
	db, err := znakdb.New(info)
	if err != nil {
		mu.mutex.Unlock()
		return nil, fmt.Errorf("repo lock open %v error %w", dbscan.TrueZnak, err)
	}
	return db, nil
}

func (r *Repository) UnlockZnak(db *znakdb.DbZnak) error {
	if db == nil {
		return fmt.Errorf("%s unlock db %v is nil", modError, dbscan.TrueZnak)
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
