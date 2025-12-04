package innfias

import (
	"sync/atomic"
)

var reentranceSearchFlag int64

func (t *page) Search() {
	defer func() {
		if r := recover(); r != nil {
			t.Logger().Error("innfias search panic %v", r)
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
}
