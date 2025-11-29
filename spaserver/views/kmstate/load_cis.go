package kmstate

import (
	"fmt"
	"korrectkm/domain/cis"

	"github.com/mechiko/utility"
)

// если csv то без экранирования!
func (m *KmStateModel) loadCisFromFile() error {
	if m.File == "" {
		return fmt.Errorf("file name empty")
	}
	if !utility.PathOrFileExists(m.File) {
		return fmt.Errorf("file %s not found", m.File)
	}
	cises, err := utility.ReadTextStringArray(m.File)
	if err != nil {
		return fmt.Errorf("file read csv error %w", err)
	}
	m.CisIn = make([]string, 0, len(cises))
	for iRow, row := range cises {
		if len(row) < 1 {
			return fmt.Errorf("row %d count column in file must be greate than zero", iRow)
		}
		cis, err := cis.TrimCis(row)
		if err != nil {
			return fmt.Errorf("row %d error %w", iRow, err)
		}
		m.CisIn = append(m.CisIn, cis)
	}
	return nil
}
