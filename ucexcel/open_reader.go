package ucexcel

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

func (ue *ucexcel) Open() error {
	var reader io.Reader
	switch ue.template {
	case "utilizedReport":
		// reader = bytes.NewReader(utilizedReportExcel)
	case "utilizedReportGtin":
		// reader = bytes.NewReader(utilizedReportGtinExcel)
	case "":
		// создаем чистый файл без стилей
		ue.file = excelize.NewFile()
		return nil
	default:
		err := fmt.Errorf("ucexcel wrong template name")
		return err
	}
	if file, err := excelize.OpenReader(reader); err != nil {
		return fmt.Errorf("%w", err)
	} else {
		ue.file = file
		ue.celStyleDefault, _ = ue.file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "right",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
			Font: &excelize.Font{
				Bold:   false,
				Family: "Arial",
				Size:   9,
			},
		})
		ue.celStyleLeft, _ = ue.file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
			Font: &excelize.Font{
				Bold:   false,
				Family: "Arial",
				Size:   9,
			},
		})
		ue.celStyleBold, _ = ue.file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
			Font: &excelize.Font{
				Bold:   true,
				Family: "Arial",
				Size:   9,
			},
		})
		ue.celStyleBoldRight, _ = ue.file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "right",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
			Font: &excelize.Font{
				Bold:   true,
				Family: "Arial",
				Size:   9,
			},
		})
		ue.celStyleVCenter, _ = ue.file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
			Font: &excelize.Font{
				Bold:   false,
				Family: "Arial",
				Size:   9,
			},
		})
		ue.celStyleBoldCenter, _ = file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
			Font: &excelize.Font{
				Bold:   true,
				Family: "Arial",
				Size:   9,
			},
		})
		ue.celStyleBoldGreen, _ = ue.file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
			Font: &excelize.Font{
				Bold:   true,
				Family: "Arial",
				Size:   9,
				Color:  "1e685a",
			},
		})
		ue.celStyleBoldBlue, _ = ue.file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
			Font: &excelize.Font{
				Bold:   true,
				Family: "Arial",
				Size:   9,
				Color:  "056c91",
			},
		})
	}

	return nil
}
