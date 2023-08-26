package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/psxzz/backend-trainee-assignment/internal/app/endpoint"
	"github.com/psxzz/backend-trainee-assignment/internal/app/service"
	"github.com/psxzz/backend-trainee-assignment/internal/app/storage/postgresql"
	"github.com/psxzz/backend-trainee-assignment/internal/config"
)

type App struct {
	db   *sql.DB
	cfg  *config.Config
	svc  *service.Service
	endp *endpoint.Endpoint
	echo *echo.Echo
}

func New() (*App, error) {
	var err error
	cfg := config.New()
	app := &App{cfg: cfg}

	db, err := sql.Open("postgres", app.cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("invalid database connection credentials: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("couldn't connect to database: %w", err)
	}

	storage := postgresql.New(db)
	app.svc, err = service.New(storage)
	if err != nil {
		return nil, fmt.Errorf("couldn't create a service: %w", err)
	}

	app.endp = endpoint.New(app.svc)

	app.echo = echo.New()

	// TODO: Declare endpoint handlers
	app.echo.GET("/", handler)

	return app, nil
}

func (a App) Run() {
	go func() {
		err := a.echo.Start(":8080")
		if err != nil {
			log.Fatal(err)
		}
	}()
	defer a.db.Close()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := a.echo.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello world!")
}
