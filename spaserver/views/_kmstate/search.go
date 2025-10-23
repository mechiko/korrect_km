package stats

import (
	"firstwails/domain"
	"firstwails/utility"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

var reentranceSearchFlag int64

func (t *page) Search(model domain.Model) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errStr := fmt.Sprintf("%s panic %v", modError, r)
				t.Logger().Error(errStr)
				modelRecover := t.Reductor().Model()
				modelRecover.Stats.Errors = append(modelRecover.Stats.Errors, errStr)
				t.UpdateModel(modelRecover, "stats.search.recover")
			}
		}()
		// эта часть в фоне выполняется только в одном экземпляре
		if atomic.CompareAndSwapInt64(&reentranceSearchFlag, 0, 1) {
			defer atomic.StoreInt64(&reentranceSearchFlag, 0)
		} else {
			t.Logger().Errorf("%s reenter page:stats:search", modError)
			return
		}
		// создаем клиента и если надо авторизуемся
		// там же сохраняется в конфиг если происходит обновление токенов
		tc := t.StartTrueClient(&model)
		if len(tc.Errors()) > 0 {
			model.Stats.Errors = append(model.Stats.Errors, tc.Errors()...)
			t.Logger().Debugf("%s trueclient authorised errors %d", modError, len(tc.Errors()))
			return
		}
		t.Logger().Debugf("%s trueclient authorised", modError)
		// после авторизации обновляем модель

		chunkSize := 1000
		chunks := utility.SplitStringSlice2Chunks(model.Stats.CisIn, chunkSize)
		t.Logger().Debugf("%s len allCises %d chunks %d", modError, len(model.Stats.CisIn), len(chunks))
		if err := t.Repo().DbLite().CisRequestDeleteAll(); err != nil {
			model.Stats.Errors = append(model.Stats.Errors, tc.Errors()...)
			t.Logger().Debugf("%s trueclient CisRequestDeleteAll errors %d", modError, len(tc.Errors()))
		}

		cisStatus := make(map[string]map[string]int)
		cisOut := make(domain.CisSlice, 0)
		start := time.Now()
		for i, chunk := range chunks {
			progress := i * 100 / len(chunks)
			t.Logger().Debugf("%s search progress %d", modError, progress)
			model.Stats.Progress = progress
			t.UpdateModel(model, "stats.search")
			cisResponce := []domain.CisPostJson{}
			if _, err := tc.CisesListPost(&cisResponce, chunk); err != nil {
				model.Stats.Errors = append(model.Stats.Errors, err.Error())
				t.Logger().Debugf("%s tc.CisesList %s", modError, err.Error())
				// при ошибке весь chunk делаем ответом с ошибкой
				for i := range chunk {
					resp := domain.CisPostJson{}
					resp.Result.Cis = chunk[i]
					resp.Result.Status = err.Error()
					cisResponce = append(cisResponce, resp)
				}
			}
			// вставка данных в бд
			if err := t.Repo().DbLite().InsertCisRequestPost(cisResponce); err != nil {
				model.Stats.Errors = append(model.Stats.Errors, tc.Errors()...)
				t.Logger().Debugf("%s InsertCisRequest %s", modError, err.Error())
			}
			for _, cisItem := range cisResponce {
				// status := domain.DictCisTypes(cisItem.Result.Status)
				status := strings.ToLower(cisItem.Result.Status)
				statusEx := strings.ToLower(cisItem.Result.StatusEx)
				if _, ok := cisStatus[status]; !ok {
					cisStatus[status] = make(map[string]int)
					cisStatus[status][statusEx] = 1
				} else {
					if count, ok := cisStatus[status][statusEx]; !ok {
						cisStatus[status][statusEx] = 1
					} else {
						cisStatus[status][statusEx] = count + 1
					}
				}
				cisOut = append(cisOut, &domain.Cis{Cis: cisItem.Result.Cis, Status: status, StatusEx: statusEx})
			}
		}
		t.Logger().Debugf("%s CisRequest time since %v", modError, time.Since(start))
		t.Logger().Debugf("%s CisList %d", modError, len(cisOut))
		model.Stats.CisOut = cisOut
		model.Stats.CisStatus = cisStatus
		model.Stats.State = 0
		// TODO обновление модели затянуто ...
		t.UpdatePage(model, "stats", "stats.search")
	}()
}
