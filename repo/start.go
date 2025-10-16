package repo

import (
	"context"
	"fmt"
)

func (r *Repository) Run(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	// ожидаем завершения контекста
	<-ctx.Done()
	r.Logger().Infof("завершаем работу репозитория по контексту")
	return nil
}
