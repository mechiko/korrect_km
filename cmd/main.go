package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"korrectkm/app"
	"korrectkm/checkdbg"
	"korrectkm/config"
	"korrectkm/embedded"
	"korrectkm/reductor"
	"korrectkm/repo"
	"korrectkm/spaserver"
	"korrectkm/trueclient"
	"korrectkm/trueclient/mystore"
	"korrectkm/zaplog"
	"os"
	"path/filepath"

	"github.com/mechiko/dbscan"
	"github.com/mechiko/utility"
	"go.uber.org/zap"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"golang.org/x/sync/errgroup"
)

const modError = "main"

// var version = "0.0.0"
var fileExe string
var dir string

// если local true то папка создается локально
var local = flag.Bool("local", false, "")

func errMessageExit(loger *zap.SugaredLogger, title string, err error) {
	if loger != nil {
		loger.Errorf("%s %v", title, err)
	}
	utility.MessageBox(title, err.Error())
	os.Exit(-1)
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, groupCtx := errgroup.WithContext(ctx)

	cfg, err := config.New("", !*local)
	if err != nil {
		errMessageExit(nil, "ошибка конфигурации", err)
	}

	var logsOutConfig = map[string][]string{
		"logger":   {"stdout", filepath.Join(cfg.LogPath(), config.Name)},
		"echo":     {filepath.Join(cfg.LogPath(), "echo")},
		"reductor": {filepath.Join(cfg.LogPath(), "reductor")},
		"true":     {filepath.Join(cfg.LogPath(), "true")},
	}
	zl, err := zaplog.New(logsOutConfig, true)
	if err != nil {
		errMessageExit(nil, "ошибка создания логера", err)
	}

	lg, err := zl.GetLogger("logger")
	if err != nil {
		errMessageExit(nil, "ошибка получения логера", err)
	}
	loger := lg.Sugar()
	loger.Debug("zaplog started")
	loger.Infof("mode = %s", config.Mode)
	if cfg.Warning() != "" {
		loger.Infof("pkg:config warning %s", cfg.Warning())
	}
	errProcessExit := func(title string, err error) {
		errMessageExit(loger, title, err)
	}

	reductorLogger, err := zl.GetLogger("reductor")
	if err != nil {
		errProcessExit("Ошибка получения логера для редуктора", err)
	}

	if err := reductor.New(reductorLogger.Sugar()); err != nil {
		errProcessExit("Ошибка создания редуктора", err)
	}

	loger.Info("new webapp")
	// создаем приложение с опциями из конфига и логером основным
	app := app.New(cfg, loger, dir)
	app.SetDbSelfPath(cfg.ConfigPath())
	// бд основные находятся в текущем каталоге если не переопределено в настройках
	app.SetDefaultDbPath("")

	// инициализируем пути необходимые приложению
	app.CreatePath()

	loger.Info("start repo")
	// инициализируем REPO

	listDbs := make(dbscan.ListDbInfoForScan)
	listDbs[dbscan.Config] = &dbscan.DbInfo{}
	listDbs[dbscan.Other] = &dbscan.DbInfo{
		File:   "korrectkm.db",
		Name:   "korrectkm",
		Driver: "sqlite",
		Path:   app.DbSelfPath(),
	}
	listDbs[dbscan.TrueZnak] = &dbscan.DbInfo{}

	err = repo.New(listDbs, ".")
	if err != nil {
		utility.MessageBox("Ошибка запуска репозитория", err.Error())
		os.Exit(-1)
	}
	repoStart, err := repo.GetRepository()
	if err != nil {
		utility.MessageBox("Ошибка получения репозитория", err.Error())
		os.Exit(-1)
	}

	// создаем редуктор с новой моделью
	modelTcl := trueclient.TrueClientModel{}
	// читаем модель из файла toml
	modelTcl.Read(app)
	// выставляем присутствие базы конфиг.дб
	modelTcl.IsConfigDB = repoStart.Is(dbscan.Config)
	// если настройка использовать авторизацию алкохелпа то загружаем данные из config.db
	if modelTcl.UseConfigDB {
		if repoStart.Is(dbscan.Config) {
			dbCfg, err := repoStart.LockConfig()
			if err != nil {
				utility.MessageBox("Ошибка получения конфигурации config.db", err.Error())
				os.Exit(-1)
			}
			defer repoStart.UnlockConfig(dbCfg)
			modelTcl.OmsID, err = dbCfg.Key("oms_id")
			if err != nil {
				utility.MessageBox("Ошибка получения конфигурации config.db", err.Error())
				os.Exit(-1)
			}
			modelTcl.DeviceID, err = dbCfg.Key("connection_id")
			if err != nil {
				utility.MessageBox("Ошибка получения конфигурации config.db", err.Error())
				os.Exit(-1)
			}
			modelTcl.HashKey, err = dbCfg.Key("certificate_thumbprint")
			if err != nil {
				utility.MessageBox("Ошибка получения конфигурации config.db", err.Error())
				os.Exit(-1)
			}
			modelTcl.TokenSUZ, err = dbCfg.Key("token_suz")
			if err != nil {
				utility.MessageBox("Ошибка получения конфигурации config.db", err.Error())
				os.Exit(-1)
			}
			modelTcl.TokenGIS, err = dbCfg.Key("token_gis_mt")
			if err != nil {
				utility.MessageBox("Ошибка получения конфигурации config.db", err.Error())
				os.Exit(-1)
			}
		}
	}
	// загружаем сертификаты пользователя
	if modelTcl.MyStore, err = mystore.List(loger); err != nil {
		loger.Errorf("%s", err.Error())
	}

	reductor.Instance().SetModel(modelTcl, false)

	group.Go(func() error {
		go func() {
			<-groupCtx.Done()
			repoStart.Shutdown()
		}()
		return repoStart.Run(groupCtx)
	})
	// тесты
	if err := checkdbg.NewChecks(webApp).Run(); err != nil {
		loger.Errorf("check error %v", err)
		cancel()
	}

	loger.Info("start up webapp")

	port := cfg.Configuration().HostPort
	if port == "" || port == "auto" {
		if portFree, err := utility.GetFreePort(); err == nil {
			port = fmt.Sprintf("%d", portFree)
			// порт не записываем в файл конфигурации остается только в модели приложения
			webApp.Config().SetInConfig("hostport", port, false)
		}
	}
	loger.Infof("http port %s", port)

	// тут инициализируются так же модели для всех видов
	httpServer := spaserver.New(webApp, port, true)
	// запускаем сервер эхо через него SSE работает для флэш сообщений
	httpServer.Start()
	// для отладки посмотреть редуктор после инициализации
	// rdct := reductor.Instance()
	// loger.Info("start wails %v", rdct.Model(reductor.Home))

	if err := httpServer.PingSetup(); err != nil {
		httpServer.SetFlush(err.Error(), "error")
		httpServer.SetActivePage(reductor.Setup)
		loger.Errorf("%s", err.Error())
	}
	if err := wails.Run(&options.App{
		Title:     "Утилиты для ЧЗ и А3",
		Width:     1040,
		Height:    768,
		MinWidth:  200,
		MinHeight: 200,
		// MaxWidth:      1280,
		// MaxHeight:     800,
		DisableResize: false,
		Fullscreen:    false,
		Frameless:     false,
		// CSSDragProperty:   "widows",
		// CSSDragValue:      "1",
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		AssetServer: &assetserver.Options{
			Assets: embedded.Root,
			// Middleware: func(next http.Handler) http.Handler {
			// 	// устанавливаем обработку not found на предлагаемую по умолчанию wails
			// 	// это произойдет когда наш роутер не найдет нужного
			// 	httpServer.Echo().RouteNotFound("/", func(c echo.Context) error {
			// 		// return c.NoContent(204)
			// 		return c.String(200, "not found")
			// 	})
			// 	return httpServer.Handler()
			// },
			Handler: httpServer.Handler(),
		},
		// Menu:             webApp.ApplicationMenu(),
		EnableDefaultContextMenu: true,
		Logger:                   nil,
		LogLevel:                 logger.INFO,
		OnStartup:                webApp.Startup,
		// OnDomReady:               httpServer.Publish,
		OnBeforeClose:    webApp.BeforeClose,
		OnShutdown:       webApp.Shutdown,
		WindowStartState: options.Normal,
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
			ZoomFactor:          1.0,
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
	}); err != nil {
		loger.Errorf("%s wails error %s", modError, err.Error())
	}
	cancel()
	// ожидание завершения всех в группе
	if err := group.Wait(); err != nil {
		fmt.Printf("game over! error %s\n", err.Error())
	} else {
		fmt.Println("game over!")
	}
	zaplog.OnShutdown()
}
