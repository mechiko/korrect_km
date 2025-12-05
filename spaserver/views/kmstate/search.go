package kmstate

import (
	"korrectkm/domain"
	"korrectkm/domain/models/modeltrueclient"
	"korrectkm/reductor"
	"korrectkm/trueclient"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/mechiko/utility"
)

var reentranceSearchFlag int64

func (t *page) Search() {
	defer func() {
		if r := recover(); r != nil {
			t.Logger().Error("kmstate search panic %v", r)
		}
	}()
	// эта часть в фоне выполняется только в одном экземпляре
	if atomic.CompareAndSwapInt64(&reentranceSearchFlag, 0, 1) {
		// defer atomic.StoreInt64(&reentranceSearchFlag, 0)
		defer func() {
			atomic.StoreInt64(&reentranceSearchFlag, 0)
		}()
	} else {
		t.Logger().Errorf("reenter page:kmstate:search")
		return
	}
	data, err := t.PageModel()
	if err != nil {
		t.Logger().Errorf("page:kmstate:search model %v", err)
		utility.MessageBox("ошибка программы", err.Error())
		return
	}
	data.Errors = make([]string, 0)
	defer func() {
		data.IsProgress = false
		data.Progress = 0
		if err := t.ModelUpdate(data); err != nil {
			t.Logger().Errorf("page:kmstate:search model %v", err)
			utility.MessageBox("ошибка программы", err.Error())
		}
	}()
	mtcl, err := reductor.Model[*modeltrueclient.TrueClientModel](domain.TrueClient)
	if err != nil {
		t.Logger().Errorf("page:kmstate:search reductor.Model %v", err)
		data.Errors = append(data.Errors, err.Error())
		return
	}
	tc, err := trueclient.NewFromModelSingle(t, mtcl)
	if err != nil {
		t.Logger().Errorf("page:kmstate:search model %v", err)
		data.Errors = append(data.Errors, err.Error())
		return
	}
	chunkSize := 1000
	chunks := slices.Chunk(data.CisIn, chunkSize)
	cisStatus := make(map[string]map[string]int)
	cisOut := make(domain.CisSlice, 0)
	// start := time.Now()
	deltaProgressChunk := float64(100)
	if len(data.CisIn) > chunkSize {
		numChunks := (len(data.CisIn) + chunkSize - 1) / chunkSize
		deltaProgressChunk = 100.0 / float64(numChunks)
	}
	progress := 0.0
	for chunk := range chunks {
		cisResponce := []domain.CisPostJson{}
		if _, err := tc.CisesListPost(&cisResponce, chunk); err != nil {
			data.Errors = append(data.Errors, err.Error())
			t.Logger().Debugf("tc.CisesList %v", err)
			// при ошибке весь chunk делаем ответом с ошибкой
			for i := range chunk {
				resp := domain.CisPostJson{}
				resp.Result.Cis = chunk[i]
				resp.Result.Status = err.Error()
				cisResponce = append(cisResponce, resp)
			}
		}
		// 	// вставка данных в бд
		// 	if err := t.Repo().DbLite().InsertCisRequestPost(cisResponce); err != nil {
		// 		model.Stats.Errors = append(model.Stats.Errors, tc.Errors()...)
		// 		t.Logger().Debugf("%s InsertCisRequest %s", modError, err.Error())
		// 	}
		for _, cisItem := range cisResponce {
			// status := domain.DictCisTypes(cisItem.Result.Status)
			status := strings.ToLower(cisItem.Result.Status)
			statusEx := strings.ToLower(cisItem.Result.StatusEx)
			if cisItem.ErrorMessage != "" {
				status = strings.ToLower(cisItem.ErrorMessage)
			}
			if cisItem.ErrorCode != "" {
				statusEx = strings.ToLower(cisItem.ErrorCode)
			}
			if cisStatus[status] == nil {
				cisStatus[status] = make(map[string]int)
			}
			cisStatus[status][statusEx]++
			cisOut = append(cisOut, &domain.Cis{Cis: cisItem.Result.Cis, Status: status, StatusEx: statusEx})
		}
		progress += deltaProgressChunk
		data.Progress = int(progress)
		reductor.SetModel(data, false)
	}
	// data, _ = t.PageModel()
	// t.Logger().Debugf("%s CisRequest time since %v", modError, time.Since(start))
	// t.Logger().Debugf("%s CisList %d", modError, len(cisOut))
	data.CisOut = cisOut
	data.CisStatus = cisStatus
	// по завершению прописываем в модель прогресс завершился
	// через defer() выше
}
