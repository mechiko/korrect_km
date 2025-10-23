package ucexcel

import (
	"fmt"
	"korrectkm/ucexcel/address"
)

func (ue *ucexcel) ColumnReport(report []string) (err error) {
	// https://xuri.me/excelize/ru/workbook.html#GetSheetProps
	ue.sheet = "Sheet1"
	countRow = 1
	ue.address = address.New(1, 0)
	for _, ss := range report {
		ue.templateReportLineColumn(ss)
	}
	return nil
}

func (ue *ucexcel) templateReportLineColumn(s string) error {
	if err := ue.file.SetCellStr(ue.sheet, ue.address.Address(), s); err != nil {
		return fmt.Errorf("excel error %w", err)
	}
	ue.address.NextRow()
	return nil
}
