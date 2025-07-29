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
	"korrectkm/utility"
	"korrectkm/zaplog"
	"os"
	"path/filepath"
	"strings"

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

func init() {
	flag.Parse()
	fileExe = os.Args[0]
	dir, _ = filepath.Abs(filepath.Dir(fileExe))
	os.Chdir(dir)
	if *local {
		if filepath.IsAbs(config.ConfigPath) {
			utility.AbsPathCreate(config.ConfigPath)
		} else {
			utility.PathCreate(config.ConfigPath)
		}
		if filepath.IsAbs(config.LogPath) {
			utility.AbsPathCreate(config.LogPath)
		} else {
			utility.PathCreate(config.LogPath)
		}
		if filepath.IsAbs(config.DbPath) {
			utility.AbsPathCreate(config.DbPath)
		} else {
			utility.PathCreate(config.DbPath)
		}
		zaplog.Run(config.LogPath, config.Name)
	} else {
		// создаем папки по конфигурации
		if filepath.IsAbs(config.ConfigPath) {
			utility.AbsPathCreate(config.ConfigPath)
		} else {
			utility.HomePathCreate(config.ConfigPath)
		}
		if filepath.IsAbs(config.LogPath) {
			utility.AbsPathCreate(config.LogPath)
		} else {
			utility.HomePathCreate(config.LogPath)
		}
		if filepath.IsAbs(config.DbPath) {
			utility.AbsPathCreate(config.DbPath)
		} else {
			utility.HomePathCreate(config.DbPath)
		}
		zaplog.Run(filepath.Join(utility.UserHomeDir(), config.LogPath), config.Name)
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, groupCtx := errgroup.WithContext(ctx)

	loger := zaplog.Logger
	loger.Debug("zaplog started")
	loger.Infof("mode = %s", config.Mode)

	cfg, err := config.New("", !*local)
	if err != nil {
		loger.Errorf("%s %s", modError, err.Error())
		panic(err.Error())
	}
	// создаем папку вывода после инициализации конфигурации
	if *local {
		if filepath.IsAbs(cfg.Configuration().Output) {
			utility.AbsPathCreate(cfg.Configuration().Output)
		} else {
			utility.PathCreate(cfg.Configuration().Output)
		}
	} else {
		if filepath.IsAbs(cfg.Configuration().Output) {
			utility.AbsPathCreate(cfg.Configuration().Output)
		} else {
			utility.HomePathCreate(filepath.Join(config.LogPath, cfg.Configuration().Output))
		}
	}

	loger.Infof("путь шаблонов %s", config.RootPathTemplates())

	// создаем редуктор для хранения моделей приложения
	reductor.New(zaplog.Reductor.Sugar())

	loger.Info("new webapp")
	// инитим роутер для http, конфиг и прочее
	webApp := app.NewWebApp(cfg, zaplog.Logger, dir)

	loger.Info("start repo")
	// инициализируем REPO

	dbPath := config.DbPath
	if !filepath.IsAbs(dbPath) {
		dbPath = filepath.Join(config.UserHomeDir, dbPath)
	}
	repoStart := repo.New(webApp, dbPath)
	if len(repoStart.Errors()) > 0 {
		fullErr := strings.Join(repoStart.Errors(), "\n")
		utility.MessageBox("Ошибки запуска репозитория", fullErr)
		os.Exit(-1)
	}
	webApp.SetRepo(repoStart)

	// создаем редуктор с новой моделью
	modelTcl := trueclient.TrueClientModel{}
	// читаем модель из файла toml
	modelTcl.Read(webApp)
	// выставляем присутствие базы конфиг.дб
	modelTcl.IsConfigDB = repoStart.IsConfig()
	// если настройка использовать авторизацию алкохелпа то загружаем данные из config.db
	if modelTcl.UseConfigDB {
		if repoStart.IsConfig() {
			modelTcl.OmsID = repoStart.ConfigDB().Key("oms_id")
			modelTcl.DeviceID = repoStart.ConfigDB().Key("connection_id")
			modelTcl.HashKey = repoStart.ConfigDB().Key("certificate_thumbprint")
			modelTcl.TokenSUZ = repoStart.ConfigDB().Key("token_suz")
			modelTcl.TokenGIS = repoStart.ConfigDB().Key("token_gis_mt")
		}
	}
	// загружаем сертификаты пользователя
	if modelTcl.MyStore, err = mystore.List(loger); err != nil {
		loger.Errorf("%s", err.Error())
	}

	reductor.Instance().SetModel(reductor.TrueClient, modelTcl)

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
