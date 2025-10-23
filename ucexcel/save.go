package ucexcel

import (
	"fmt"
	"os"
	_ "time/tzdata"
)

func (ue *ucexcel) Save(fileAbs string) error {
	file, err := os.OpenFile(fileAbs, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer func() error {
		if err := file.Close(); err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	}()

	if _, err := ue.file.WriteTo(file); err != nil {
		return fmt.Errorf("%w", err)
	}

	return err
}
