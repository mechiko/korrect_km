package guiconnect

import (
	"fmt"
	"korrectkm/domain"
	"korrectkm/domain/models/modeltrueclient"
	"time"

	"github.com/mechiko/utility"
	"github.com/mechiko/walk"
	dcl "github.com/mechiko/walk/declarative"
)

type StoreItem struct {
	Hash string
	Name string
}

func StartDialog(app domain.Apper, model *modeltrueclient.TrueClientModel) error {
	var dlg *walk.Dialog
	var acceptPB, cancelPB, pingPB *walk.PushButton
	var useConfigDB *walk.CheckBox
	var storeCB *walk.ComboBox
	var DB *walk.DataBinder
	// var userClose bool
	var omsID, deviceID *walk.TextEdit

	mdlStore := make([]*StoreItem, 0)
	for h, n := range model.MyStore {
		mdlStore = append(mdlStore, &StoreItem{Hash: h, Name: n})
	}

	icon, err := walk.Resources.Icon("3")
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := (dcl.Dialog{
		AssignTo:      &dlg,
		Title:         "Подключение к ЧЗ",
		Size:          dcl.Size{Width: 400, Height: 400},
		Icon:          icon,
		Layout:        dcl.VBox{Spacing: 10, Margins: dcl.Margins{Left: 5, Top: 5, Right: 5, Bottom: 5}},
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: dcl.DataBinder{
			AssignTo:        &DB,
			DataSource:      model,
			AutoSubmit:      true,
			AutoSubmitDelay: 1 * time.Second,
		},
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
						Checked:        dcl.Bind("UseConfigDB"),
						Enabled:        dcl.Bind("IsConfigDB"),
						OnCheckedChanged: func() {
							if useConfigDB.Checked() {
								model.UseConfigDB = true
								app.SetOptions("trueclient.useconfigdb", true)
								err := model.ReadState(app)
								if err != nil {
									utility.MessageBox("ошибка чтения модели", err.Error())
								}
								deviceID.SetText(model.DeviceID)
								omsID.SetText(model.OmsID)
								si := &StoreItem{Hash: model.HashKey, Name: model.MyStore[model.HashKey]}
								idx := -1
								for i, v := range mdlStore {
									if (v.Hash == si.Hash) && (v.Name == si.Name) {
										idx = i
										break // Exit the loop once found
									}
								}
								// idx := slices.Index(mdlStore, si)
								storeCB.SetCurrentIndex(idx)
								storeCB.SetEnabled(false)
								deviceID.SetEnabled(false)
								omsID.SetEnabled(false)
							} else {
								model.UseConfigDB = false
								app.SetOptions("trueclient.useconfigdb", false)
								deviceID.SetText("")
								storeCB.SetEnabled(true)
								deviceID.SetEnabled(true)
								omsID.SetEnabled(true)
							}
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
								Text: "КЭП:",
							},
							dcl.ComboBox{
								AssignTo:      &storeCB,
								Name:          "combobox",
								BindingMember: "Hash",
								DisplayMember: "Name",
								Editable:      false,
								Value:         dcl.Bind("HashKey"),
								Model:         mdlStore,
								OnCurrentIndexChanged: func() {
									// idx := storeCB.CurrentIndex()
									// fmt.Println("combo - ", idx, " ", mdlStore[idx])
								},
							},
						},
					},
					dcl.Composite{
						Layout:    dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
						Border:    false,
						Alignment: dcl.Alignment2D(walk.AlignHNearVNear),
						Children: []dcl.Widget{
							dcl.Label{
								Text: "OMS ID:",
							},
							dcl.TextEdit{
								AssignTo: &omsID,
								Text:     dcl.Bind("OmsID"),
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
								Text: "Device ID:",
							},
							dcl.TextEdit{
								AssignTo: &deviceID,
								Text:     dcl.Bind("DeviceID"),
							},
							dcl.HSpacer{},
						},
					},
				}},
			dcl.Composite{
				Border: false,
				Layout: dcl.HBox{MarginsZero: true, SpacingZero: true, Margins: dcl.Margins{Left: 5, Top: 5, Right: 5, Bottom: 0}},
				Children: []dcl.Widget{
					dcl.PushButton{
						AssignTo: &pingPB,
						Text:     "Проверка соединения",
						OnClicked: func() {
							if !useConfigDB.Checked() {
								model.OmsID = omsID.Text()
								model.DeviceID = deviceID.Text()
								model.TokenGIS = ""
								model.TokenSUZ = ""
								idx := storeCB.CurrentIndex()
								if idx < 0 {
									utility.MessageBox("ошибка", "индекс КЭП меньше 0")
									return
								}
								model.HashKey = mdlStore[idx].Hash
							}
							err := ping(app, model)
							if err != nil {
								utility.MessageBox("ошибка", err.Error())
								return
							} else {
								if model.PingSuz != nil {
									utility.MessageBox("Проверка успешна", fmt.Sprintf("api %v", model.PingSuz.ApiVersion))
								} else {
									utility.MessageBox("Ошибка", "пинг отсутствует")
								}
							}
						},
					},
					dcl.PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if !useConfigDB.Checked() {
								model.OmsID = omsID.Text()
								model.DeviceID = deviceID.Text()
								idx := storeCB.CurrentIndex()
								if idx < 0 {
									utility.MessageBox("ошибка", "индекс КЭП меньше 0")
									return
								}
								model.HashKey = mdlStore[idx].Hash
							}
							err := ping(app, model)
							if err != nil {
								utility.MessageBox("ошибка установки соединения", err.Error())
								return
							}
							err = model.SyncToStore(app)
							if err != nil {
								utility.MessageBox("ошибка сохранения модели", err.Error())
								return
							}
							err = app.SaveOptions()
							if err != nil {
								utility.MessageBox("ошибка сохранения модели в файл", err.Error())
								return
							}
							dlg.Accept()
						},
					},
					dcl.PushButton{
						AssignTo:  &cancelPB,
						Text:      "Выход",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
			dcl.VSpacer{},
		},
	}).Create(nil); err != nil {
		return fmt.Errorf("%w", err)
	}
	// This should do it (or not :-p)
	// dlg.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
	// 	// если закрывается окно диалога крестом
	// 	if reason == walk.CloseReasonUnknown {
	// 		userClose = true
	// 	}
	// })

	if model.UseConfigDB && model.IsConfigDB {
		storeCB.SetEnabled(false)
		deviceID.SetEnabled(false)
		omsID.SetEnabled(false)
	}
	if ret := dlg.Run(); ret != 1 {
		return fmt.Errorf("dialog return %d", ret)
	}
	// if userClose {
	// 	return fmt.Errorf("dialog user close")
	// }
	return nil
}
