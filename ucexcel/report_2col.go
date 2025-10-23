package ucexcel

import (
	"fmt"
	"korrectkm/ucexcel/address"
)

var countRow int = 1

func (ue *ucexcel) TwoColumnReport(report [][]string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panic %v", modError, r)
		}
	}()
	// https://xuri.me/excelize/ru/workbook.html#GetSheetProps
	ue.sheet = "Sheet1"
	countRow = 1
	ue.address = address.New(1, 0)
	for _, ss := range report {
		ue.templateReportListLine(ss)
	}
	return nil
}

func (ue *ucexcel) templateReportListLine(s []string) error {
	if err := ue.file.SetCellStr(ue.sheet, ue.address.Address(), s[0]); err != nil {
		return fmt.Errorf("excel error %w", err)
	}
	ue.address.NextCol()
	if err := ue.file.SetCellStr(ue.sheet, ue.address.Address(), s[1]); err != nil {
		return fmt.Errorf("excel error %w", err)
	}
	ue.address.NextRow()
	return nil
}
