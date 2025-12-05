package ucexcel

import (
	"fmt"
	"time"
	_ "time/tzdata"
)

const ext = "xlsx"

// func (ue *ucexcel) ExcelFileName(prefix string) string {
// 	startDate := fmt.Sprintf("%02d.%02d.%4d", ue.StartDate().Local().Day(), ue.StartDate().Local().Month(), ue.StartDate().Local().Year())
// 	endDate := fmt.Sprintf("%02d.%02d.%4d", ue.EndDate().Local().Day(), ue.EndDate().Local().Month(), ue.EndDate().Local().Year())
// 	// guid := utility.String(8)
// 	// name := prefix + " " + startDate + " " + endDate + " " + guid
// 	// name += "." + ext
// 	name := fmt.Sprintf("%s %s %s.%s", prefix, startDate, endDate, ext)
// 	return filepath.Join("output", name)
// }

// func (ue *ucexcel) ExcelFileNameDownload(prefix string) string {
// 	startDate := fmt.Sprintf("%02d.%02d.%4d", ue.StartDate().Local().Day(), ue.StartDate().Local().Month(), ue.StartDate().Local().Year())
// 	endDate := fmt.Sprintf("%02d.%02d.%4d", ue.EndDate().Local().Day(), ue.EndDate().Local().Month(), ue.EndDate().Local().Year())
// 	name := prefix + " " + startDate + " " + endDate
// 	name += "." + ext
// 	return name
// }

func (ue *ucexcel) ExcelFileNameSimple(prefix string) string {
	ds := time.Now()
	startDate := fmt.Sprintf("%02d.%02d.%4d", ds.Local().Day(), ds.Local().Month(), ds.Local().Year())
	name := prefix + "_" + startDate
	name += "." + ext
	return name
}
