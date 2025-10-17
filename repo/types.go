package repo

import (
	"errors"
	"fmt"
	"sync"

	"github.com/mechiko/dbscan"
)

const modError = "pkg:repo"

var Version int64

type singleMutex struct {
	mutex sync.Mutex
}

var rp *Repository
var rpMu sync.RWMutex
var initMu sync.Mutex

type Repository struct {
	dbs     *dbscan.Dbs
	dbMutex map[dbscan.DbInfoType]*singleMutex
	listDbs []dbscan.DbInfoType
}

// dbPath для своей БД
// func New(logcfg ILogCfg, dbPath string) (rp *Repository, err error) {
func New(listDbs dbscan.ListDbInfoForScan, dbPath string) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = fmt.Errorf("repo panic %v", rerr)
		}
	}()
	initMu.Lock()
	defer initMu.Unlock()
	if rp != nil {
		return fmt.Errorf("repo repository already initialized")
	}
	r := &Repository{
		dbMutex: make(map[dbscan.DbInfoType]*singleMutex),
	}
	if len(listDbs) == 0 {
		return fmt.Errorf("список описателей бд пуст")
	}
	for tp, info := range listDbs {
		if info == nil {
			return fmt.Errorf("%s in list [%v] is nil", modError, tp)
		}
	}
	dbs, err := dbscan.New(listDbs, dbPath)
	if err != nil {
		return fmt.Errorf("%s dbscan error %w", modError, err)
	}
	r.dbs = dbs
	r.listDbs = make([]dbscan.DbInfoType, 0, len(listDbs))
	for tp := range listDbs {
		dbInfo := r.dbs.Info(tp)
		if dbInfo == nil {
			// такая ошибка не вероятна дбскан выдаст ошибку при сканировании
			// но проверить надо вдруг чего...
			err = errors.Join(err, fmt.Errorf("%s отсутствует БД %v", modError, tp))
		} else {
			r.listDbs = append(r.listDbs, tp)
			// создаем в мапе мьютекс
			if _, ok := r.dbMutex[tp]; !ok {
				r.dbMutex[tp] = &singleMutex{}
			}
		}
	}
	if err != nil {
		return fmt.Errorf("%s не все бд найдены %v", modError, err)
	}
	if di := r.dbs.Info(dbscan.Other); di != nil {
		// инициализация для Self если она есть в настройках списка доступных БД
		if err := r.prepareSelf(); err != nil {
			return fmt.Errorf("%s ошибка миграции self %w", modError, err)
		}
	}
	rpMu.Lock()
	rp = r
	rpMu.Unlock()
	return nil
}

func GetRepository() (*Repository, error) {
	rpMu.RLock()
	defer rpMu.RUnlock()
	if rp == nil {
		return nil, fmt.Errorf("repository not init")
	}
	return rp, nil
}

// возвращаем DbInfo или nil
func (r *Repository) Info(t dbscan.DbInfoType) *dbscan.DbInfo {
	if di := r.dbs.Info(t); di != nil {
		return di
	}
	return nil
}

func (r *Repository) ListDbs() (out []dbscan.DbInfoType) {
	if r.listDbs == nil {
		return nil
	}
	out = make([]dbscan.DbInfoType, len(r.listDbs))
	copy(out, r.listDbs)
	return out
}
