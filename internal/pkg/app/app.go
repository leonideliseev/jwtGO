package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/jwtGO/config"
	"github.com/leonideliseev/jwtGO/internal/handler"
	"github.com/leonideliseev/jwtGO/internal/repository"
	"github.com/leonideliseev/jwtGO/internal/service"
	"github.com/leonideliseev/jwtGO/pkg/postgresql"
)

type Closer interface {
	Close()
}

type App struct {
	cfg *config.Config

	srv *http.Server
	conn *pgxpool.Pool

	repo *repository.Repository
	serv *service.Service
	hand *handler.Handler

	quit chan os.Signal
}

func NewApp() (*App, error) {
	app := &App{}

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	app.cfg = cfg

	err = app.initDBConn(postgresql.Config{
		Host: cfg.Postgresql.Host,
		Port: cfg.Postgresql.Port,
		Username: cfg.Postgresql.User,
		Password: cfg.Postgresql.Password,
		DBName: cfg.Postgresql.Database,
		SSLMode: cfg.Postgresql.SSLMode,
	})
	if err != nil {
		return nil, err
	}
	app.initAppCore()
	app.initServer(cfg.HTTP)
	app.initShutdown()

	return app, nil
}

func (a *App) Run() {
	go func() {
		if err := a.srv.ListenAndServe(); err != nil {
			log.Fatalf("error running server: %s", err.Error())
		}
	}()

	log.Printf("JWT app running on %s:%s", a.cfg.HTTP.Host, a.cfg.HTTP.Port)

	<-a.quit

	if err := a.srv.Close(); err != nil {
		log.Print("error close server")
	}
	a.conn.Close()
}

func (a *App) initAppCore() {
	a.repo = repository.New(a.conn)
	a.serv = service.New(a.repo, a.cfg.JWT)
	a.hand = handler.New(a.serv)
}

func (a *App) initServer(cfg config.HTTP) {
	router := a.hand.InitRoutes()

	a.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: router,
	}
}

func (a *App) initShutdown() {
	a.quit = make(chan os.Signal, 1)
	signal.Notify(a.quit, syscall.SIGTERM, syscall.SIGINT)
}

