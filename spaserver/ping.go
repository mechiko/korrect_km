package spaserver

import (
	"fmt"
	"korrectkm/reductor"
	"korrectkm/trueclient"
)

// при запуске программы первый пинг блокирующий для проверки
func (s *Server) PingSetup() error {
	model, ok := reductor.Instance().Model(reductor.TrueClient).(trueclient.TrueClientModel)
	if !ok {
		return fmt.Errorf("объект редуктора не соответствует trueclient.TrueClientModel")
	}

	tcl, err := trueclient.NewFromModelSingle(s, model)
	if err != nil {
		return fmt.Errorf("failed to create trueclient: %w", err)
	}

	png, err := tcl.PingSuz()
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	model.PingSuz = png
	reductor.Instance().SetModel(reductor.TrueClient, model)
	return nil
}
