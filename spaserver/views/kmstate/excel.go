package kmstate

import (
	"fmt"
	"korrectkm/ucexcel"
	"slices"
)

func (t *page) ToExcel(ar []string, name string, size int) (fileName string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic excel %v", r)
		}
	}()
	// chunks := utility.SplitStringSlice2Chunks(ar, size)
	chunks := slices.Chunk(ar, size)
	i := 0
	for chunk := range chunks {
		fileName = fmt.Sprintf("%s_%0d[%0d].xlsx", name, i*size+1, len(chunk))
		i++
		excel := ucexcel.New(t, "", "", fileName)
		if err := excel.Open(); err != nil {
			return "", fmt.Errorf("%w", err)
		}
		if err := excel.ColumnReport(chunk); err != nil {
			return "", fmt.Errorf("%w", err)
		}
		if err := excel.Save(fileName); err != nil {
			return "", fmt.Errorf("%w", err)
		}
	}
	return fileName, nil
}
