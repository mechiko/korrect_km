package spaserver

import (
	"context"
	"fmt"
	"korrectkm/domain"
	"korrectkm/spaserver/templates"
	"korrectkm/spaserver/views"
	"korrectkm/sse"
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
	_defaultAddr            = "127.0.0.1:8888"
	_defaultShutdownTimeout = 1 * time.Second
)

// Server -.
type Server struct {
	domain.Apper
	addr            string
	server          *echo.Echo
	notify          chan error
	shutdownTimeout time.Duration
	sessionManager  *scs.SessionManager
	debug           bool
	private         *echo.Group
	templates       templates.ITemplateUI
	views           map[domain.Model]views.IView
	menu            []domain.Model
	activePage      domain.Model
	defaultPage     string
	flush           *FlushMsg
	flushMu         sync.RWMutex
	htmx            *htmx.HTMX
	sseManager      *sse.Server
	streamError     *sse.Stream
	streamInfo      *sse.Stream
}

// var sseManager *sse.Server

func New(a domain.Apper, eLogger *zap.Logger, port string, debug bool) *Server {
	addr := fmt.Sprintf("%s:%s", "127.0.0.1", port)
	if port == "" {
		addr = _defaultAddr
	}
	sess := scs.New()
	sess.Lifetime = 24 * time.Hour
	e := echo.New()
	e.Use(
		session.LoadAndSave(sess),
		zap4echo.Logger(eLogger),
		zap4echo.Recover(eLogger),
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
		Apper:           a,
		addr:            addr,
		server:          e,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		private:         private,
		debug:           debug,
		sessionManager:  sess,
		views:           make(map[domain.Model]views.IView), // массив видов по нему находим шаблоны для рендера
		menu:            make([]domain.Model, 0),
		defaultPage:     "",
		activePage:      domain.KMState,
		htmx:            htmx.New(),
	}

	e.Renderer = ss
	ss.templates = templates.New(ss, debug)
	ss.menu = append(ss.menu, domain.Home)
	ss.menu = append(ss.menu, domain.KMState)
	ss.menu = append(ss.menu, domain.InnFias)
	ss.menu = append(ss.menu, domain.Money)
	ss.Routes()
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

func (s *Server) SetActivePage(p domain.Model) {
	s.activePage = p
}

func (s *Server) ActivePage() domain.Model {
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

func (s *Server) Views() map[domain.Model]views.IView {
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

func (s *Server) Menu() []domain.Model {
	return s.menu
}
