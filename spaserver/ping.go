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
	mdl, err := reductor.Model[*modeltrueclient.TrueClientModel](domain.TrueClient)
	if err != nil {
		return fmt.Errorf("failed to create trueclient: %w", err)
	}

	tcl, err := trueclient.NewFromModelSingle(s, mdl)
	if err != nil {
		return fmt.Errorf("failed to create trueclient: %w", err)
	}

	png, err := tcl.PingSuz()
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	mdl.PingSuz = png
	return reductor.SetModel(mdl, false)
}
