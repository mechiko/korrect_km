package spaserver

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/domain/models/modeltrueclient"
	"korrectkm/reductor"
	"korrectkm/trueclient"
)

// при запуске программы первый пинг блокирующий для проверки
func (s *Server) PingSetup() error {
	mdl, err := reductor.Instance().Model(domain.TrueClient)
	if err != nil {
		return fmt.Errorf("failed to create trueclient: %w", err)
	}
	model, ok := mdl.(modeltrueclient.TrueClientModel)
	if !ok {
		return fmt.Errorf("объект редуктора не соответствует trueclient.TrueClientModel")
	}

	tcl, err := trueclient.NewFromModelSingle(s, &model)
	if err != nil {
		return fmt.Errorf("failed to create trueclient: %w", err)
	}

	png, err := tcl.PingSuz()
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	model.PingSuz = png
	return reductor.SetModel(&model, false)
}
