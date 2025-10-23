package stats

import (
	"firstwails/usecase/ucexcel"
	"firstwails/utility"
	"fmt"
)

func (t *page) ToExcel(ar []string, name string, size int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic gui excel %v", r)
		}
	}()
	chunks := utility.SplitStringSlice2Chunks(ar, size)
	for i, chunk := range chunks {
		nameFile := fmt.Sprintf("%s_%0d[%0d]", name, i*size+1, len(chunk))
		excel := ucexcel.New(t, "", "", nameFile)
		if err := excel.Open(); err != nil {
			return fmt.Errorf("%w", err)
		}
		if err := excel.ReportList(chunk); err != nil {
			return fmt.Errorf("%w", err)
		}
		if err := excel.SaveSimple(); err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil

}
