package setup

import (
	"korrectkm/reductor"
	"korrectkm/trueclient"
	"korrectkm/trueclient/mystore"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (t *page) Routes() error {
	t.Echo().GET("/setup", t.Index)
	t.Echo().POST("/setup/configdb", t.ConfigDB)
	t.Echo().POST("/setup/ping", t.Ping)
	t.Echo().POST("/setup/omsid", t.ValidateOmsID)
	// e.POST("/setup/deviceid", t.ValidateDeviceID)
	t.Echo().POST("/setup/save", t.Save)
	// e.GET("/setup/ready", t.Ready)
	return nil
}

func (t *page) Index(c echo.Context) error {
	t.ResetValidateData()
	model := t.PageModel()
	// синхронизируем с моделью TrueClientModel
	model.Read(t.Logger())
	reductor.Instance().SetModel(reductor.Setup, model)
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("index", t.PageData())); err != nil {
		return t.ServerError(c, err)
	}
	// сброс или установка извещения ошибки
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("alert", t.PageData())); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

// htmx запрос проверяем
func (t *page) ConfigDB(c echo.Context) error {
	useConfigDB := c.FormValue("useconfigdb")
	model := t.PageModel()
	if useConfigDB == "true" {
		model.ReadConfigDB(t.Repo())
	} else {
		model.ClearConfigDB()
		if store, err := mystore.List(t.Logger()); err != nil {
			return t.ServerError(c, err)
		} else {
			model.MyStore = store
		}
		// нельзя делать одинаковые DeviceId
		model.DeviceID = ""
	}
	model.PingSuz = nil
	model.Sync(t)
	reductor.Instance().SetModel(reductor.Setup, model)

	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("page", model)); err != nil {
		return t.ServerError(c, err)
	}
	return nil
}

func (t *page) Ping(c echo.Context) error {
	model := t.PageModel()
	model.DeviceID = c.FormValue("deviceid")
	model.OmsID = c.FormValue("omsid")
	model.HashKey = c.FormValue("hashkey")
	model.PingSuz = nil
	if !model.UseConfigDB {
		model.TokenGIS = ""
		model.TokenSUZ = ""
		model.AuthTime = time.Time{}
	}
	model.Validates = make(map[string]string)
	model.Errors = make([]string, 0)

	if err := uuid.Validate(model.DeviceID); err != nil {
		model.Errors = append(model.Errors, err.Error())
		model.Validates["deviceid"] = err.Error()
	}
	if err := uuid.Validate(model.OmsID); err != nil {
		model.Errors = append(model.Errors, err.Error())
		model.Validates["omsid"] = err.Error()
	}
	model.Sync(t)

	tclModel, err := model.ToTrueClient()
	if err != nil {
		return t.ServerError(c, err)
	}
	if tcl, err := trueclient.NewFromModelSingle(t, tclModel); err != nil {
		// ошибка
		return t.ServerError(c, err)
	} else {
		// если нет ошибок
		if png, err := tcl.PingSuz(); err != nil {
			// ошибка
			return t.ServerError(c, err)
		} else {
			// пинг успешен
			tclModel.PingSuz = png
			model.PingSuz = png
			reductor.Instance().SetModel(reductor.TrueClient, tclModel)
		}
	}

	reductor.Instance().SetModel(reductor.Setup, model)
	if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("index", model)); err != nil {
		return t.ServerError(c, err)
	}
	// сброс или установка извещения ошибки
	// if err := c.Render(http.StatusOK, t.Name(), t.RenderPageModel("alert", model)); err != nil {
	// 	t.ServerError(c, err)
	// }
	modelPing := map[string]interface{}{
		"template": "ping",
		"data":     tclModel,
	}
	if ping, err := t.RenderString("setup", modelPing); err != nil {
		return t.ServerError(c, err)
	} else {
		t.SetFlush(ping, "info")
	}
	return nil
}

// сохраняем настройки в конфиге и модели редуктора TrueClient
func (t *page) Save(c echo.Context) error {
	model := t.PageModel()
	model.DeviceID = c.FormValue("deviceid")
	model.OmsID = c.FormValue("omsid")
	model.HashKey = c.FormValue("hashkey")
	model.Sync(t)
	c.NoContent(204)
	t.SetFlush("Сохранено", "info")
	return nil
}
