package kmstate

import (
	"fmt"
	"korrectkm/ucexcel"
	"slices"
)

func (t *page) ToExcel(ar []string, name string, size int) (fileNames []string, err error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid chunk size: %d (must be > 0)", size)
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in Excel export: %v", r)
		}
	}()
	// chunks := utility.SplitStringSlice2Chunks(ar, size)
	chunks := slices.Chunk(ar, size)
	i := 0
	for chunk := range chunks {
		fileName := fmt.Sprintf("%s_%0d[%0d].xlsx", name, i*size+1, len(chunk))
		fileNames = append(fileNames, fileName)
		i++
		excel := ucexcel.New(t, "", "", fileName)
		if err := excel.Open(); err != nil {
			return nil, fmt.Errorf("failed to open Excel file %q: %w", fileName, err)
		}
		if err := excel.ColumnReport(chunk); err != nil {
			return nil, fmt.Errorf("failed to write Excel report to %q: %w", fileName, err)
		}
		if err := excel.Save(fileName); err != nil {
			return nil, fmt.Errorf("failed to save Excel file %q: %w", fileName, err)
		}
	}
	return fileNames, nil
}
