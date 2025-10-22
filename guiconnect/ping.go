package guiconnect

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/domain/models/modeltrueclient"
	"korrectkm/trueclient"
)

func ping(app domain.Apper, model *modeltrueclient.TrueClientModel) error {
	_, err := trueclient.NewFromModelSingle(app, model)
	if err != nil {
		return fmt.Errorf("failed to create trueclient: %w", err)
	}
	return nil
}
