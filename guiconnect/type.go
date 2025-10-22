package guiconnect

import (
	"fmt"

	"github.com/mechiko/walk"
	dcl "github.com/mechiko/walk/declarative"
)

func Start() error {
	var mw *walk.MainWindow
	var useConfigDB *walk.CheckBox

	icon, err := walk.Resources.Icon("3")
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if _, err := (dcl.MainWindow{
		AssignTo: &mw,
		Title:    "Подключение к ЧЗ",
		Size:     dcl.Size{Width: 400, Height: 400},
		Icon:     icon,
		Layout:   dcl.VBox{Spacing: 10, Margins: dcl.Margins{Left: 5, Top: 5, Right: 5, Bottom: 5}},
		Children: []dcl.Widget{
			dcl.Composite{
				Border:    false,
				Layout:    dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 5, Top: 5, Right: 5, Bottom: 0}},
				Alignment: dcl.Alignment2D(walk.AlignHNearVNear),
				Children: []dcl.Widget{
					dcl.CheckBox{
						AssignTo:       &useConfigDB,
						Name:           "checkUseConfigDB",
						Text:           "Использовать AlcoHelp3",
						TextOnLeftSide: true,
						Checked:        true,
						OnCheckedChanged: func() {
							fmt.Println(useConfigDB.Checked())
						},
					},
				},
			},
			dcl.GroupBox{
				Title:  "Параметры подключения",
				Layout: dcl.VBox{MarginsZero: false, SpacingZero: false, Margins: dcl.Margins{Left: 10, Top: 10, Right: 10, Bottom: 10}},
				Children: []dcl.Widget{
					dcl.Composite{
						Layout:    dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
						Border:    false,
						Alignment: dcl.Alignment2D(walk.AlignHNearVNear),
						Children: []dcl.Widget{
							dcl.Label{
								Text: "БД Config:",
							},
							dcl.Label{
								// AssignTo: &p.lblDbConfig,
								// Text:     model.DbConfigDesc,
							},
							dcl.HSpacer{},
						},
					},
					dcl.Composite{
						Layout:    dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
						Border:    false,
						Alignment: dcl.Alignment2D(walk.AlignHNearVNear),
						Children: []dcl.Widget{
							dcl.Label{
								Text: "БД A3:",
							},
							dcl.Label{
								// AssignTo: &p.lblDbA3,
								// Text:     model.DbA3Desc,
							},
							dcl.HSpacer{},
						},
					},
					dcl.Composite{
						Layout:    dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
						Border:    false,
						Alignment: dcl.Alignment2D(walk.AlignHNearVNear),
						Children: []dcl.Widget{
							dcl.Label{
								Text: "БД ЧЗ:",
							},
							dcl.Label{
								// AssignTo: &p.lblDbZnak,
								// Text:     model.DbZnakDesc,
							},
							dcl.HSpacer{},
						},
					},
					dcl.Composite{
						Layout:    dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
						Border:    false,
						Alignment: dcl.Alignment2D(walk.AlignHNearVNear),
						Children: []dcl.Widget{
							dcl.Label{
								Text: "БД программы:",
							},
							dcl.Label{
								// AssignTo: &p.lblDbLite,
								// Text:     model.DbLiteDesc,
							},
							dcl.HSpacer{},
						},
					},
				}},
			dcl.Composite{
				Border: false,
				Layout: dcl.VBox{MarginsZero: true, SpacingZero: true, Margins: dcl.Margins{Left: 5, Top: 5, Right: 5, Bottom: 0}},
				Children: []dcl.Widget{
					dcl.PushButton{
						Text: "Открыть папку выгрузки",
						// OnClicked: func() {
						// 	out := p.Options().Output
						// 	if out == "" {
						// 		utility.MessageBox("ошибка", "путь выгрузки не настроен")
						// 		return
						// 	}
						// 	if err := utility.OpenFileInShell(out); err != nil {
						// 		utility.MessageBox("ошибка", err.Error())
						// 	}
						// },
					},
					dcl.PushButton{
						Text: "Открыть папку настройки и логов",
						// OnClicked: func() {
						// 	cfg := p.ConfigPath()
						// 	if cfg == "" {
						// 		utility.MessageBox("ошибка", "путь к настройкам/логам не настроен")
						// 		return
						// 	}
						// 	if err := utility.OpenFileInShell(cfg); err != nil {
						// 		utility.MessageBox("ошибка", err.Error())
						// 	}
						// },
					},
				}},
			dcl.VSpacer{},
		},
	}.Run()); err != nil {
		return fmt.Errorf("guiconnect error %w", err)
	}
	return nil
}
