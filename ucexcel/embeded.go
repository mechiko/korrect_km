package ucexcel

import (
	_ "embed"
)

//go:embed templates/reestrReport.xlsx
var reestrReportExcel []byte

//go:embed templates/transactionReport.xlsx
var transactionReportExcel []byte
