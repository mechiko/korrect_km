package spaserver

import (
	"context"
	"fmt"
	"korrectkm/config"
	"korrectkm/reductor"
	"korrectkm/repo"
	"korrectkm/spaserver/templates"
	"korrectkm/spaserver/views"
	"korrectkm/sse"
	"korrectkm/zaplog"
	"net/http"
	"sync"
	"time"

	"github.com/alexedwards/scs/v2"
	session "github.com/canidam/echo-scs-session"
	"github.com/donseba/go-htmx"
	"github.com/karagenc/zap4echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = "127.0.0.1:8888"
	_defaultShutdownTimeout = 1 * time.Second
)

type ILogCfg interface {
	Config() config.IConfig
	Logger() *zap.SugaredLogger
	Repo() *repo.Repository
	Ctx() context.Context
}

// Server -.
type Server struct {
	ILogCfg
	addr            string
	server          *echo.Echo
	notify          chan error
	shutdownTimeout time.Duration
	sessionManager  *scs.SessionManager
	debug           bool
	private         *echo.Group
	templates       templates.ITemplateUI
	views           map[reductor.ModelType]views.IView
	menu            []reductor.ModelType
	activePage      reductor.ModelType
	defaultPage     string
	flush           *FlushMsg
	flushMu         sync.RWMutex
	htmx            *htmx.HTMX
	sseManager      *sse.Server
	streamError     *sse.Stream
	streamInfo      *sse.Stream
}

// var sseManager *sse.Server

func New(a ILogCfg, port string, debug bool) *Server {
	addr := fmt.Sprintf("%s:%s", "127.0.0.1", port)
	if port == "" {
		addr = _defaultAddr
	}
	sess := scs.New()
	sess.Lifetime = 24 * time.Hour
	e := echo.New()
	e.Use(
		session.LoadAndSave(sess),
		zap4echo.Logger(zaplog.EchoSugar),
		zap4echo.Recover(zaplog.EchoSugar),
	)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"authorization", "Content-Type"},
		AllowCredentials: true,
		AllowMethods:     []string{echo.OPTIONS, echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	// e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
	// 	HTML5:      true,
	// 	Root:       "root", // because files are located in `root` directory
	// 	Filesystem: http.FS(embeded.Root),
	// }))
	// наследует родительские middleware
	private := e.Group("/admin")
	ss := &Server{
		ILogCfg:         a,
		addr:            addr,
		server:          e,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		private:         private,
		debug:           debug,
		sessionManager:  sess,
		views:           make(map[reductor.ModelType]views.IView), // массив видов по нему находим шаблоны для рендера
		menu:            make([]reductor.ModelType, 0),
		defaultPage:     "",
		activePage:      reductor.Home,
		htmx:            htmx.New(),
	}

	e.Renderer = ss
	ss.templates = templates.New(ss, config.RootPathTemplates(), debug)
	ss.Routes()
	ss.menu = append(ss.menu, reductor.Home)
	ss.menu = append(ss.menu, reductor.Setup)
	ss.sseManager = sse.New()
	ss.streamError = ss.sseManager.CreateStream("error")
	ss.streamInfo = ss.sseManager.CreateStream("info")
	// go func() {
	// 	e.Logger.Fatal(e.Start(addr))
	// }()
	return ss
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.server.Start(s.addr)
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Handler() http.Handler {
	return s.server
}

func (s *Server) SessionManager() *scs.SessionManager {
	return s.sessionManager
}

func (s *Server) Echo() *echo.Echo {
	return s.server
}

func (s *Server) SetActivePage(p reductor.ModelType) {
	s.activePage = p
}

func (s *Server) ActivePage() reductor.ModelType {
	return s.activePage
}

// устанавливает заголовок окна используется в Render
func (s *Server) SetTitlePage(title string) {
	runtime.WindowSetTitle(s.Ctx(), title)
}

// используется в меню
// func (s *Server) ActivePageTitle(title string) string {
// 	runtime.WindowSetTitle(s.Ctx(), title)
// 	view, ok := s.views[s.activePage]
// 	if !ok {
// 		return "нет такого вида"
// 	}
// 	return view.Title()
// }

func (s *Server) Views() map[reductor.ModelType]views.IView {
	return s.views
}

func (s *Server) Reload() {
	if s.streamError != nil && s.streamError.Eventlog != nil {
		s.streamError.Eventlog.Clear()
	}
	if s.streamInfo != nil && s.streamInfo.Eventlog != nil {
		s.streamInfo.Eventlog.Clear()
	}
	if ctx := s.Ctx(); ctx != nil {
		runtime.WindowReload(ctx)
	}
}

func (s *Server) Htmx() *htmx.HTMX {
	return s.htmx
}

func (s *Server) Menu() []reductor.ModelType {
	return s.menu
}
