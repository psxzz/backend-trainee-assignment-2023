package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/psxzz/backend-trainee-assignment/internal/app/endpoint"
	"github.com/psxzz/backend-trainee-assignment/internal/app/service"
	"github.com/psxzz/backend-trainee-assignment/internal/app/storage/postgresql"
	"github.com/psxzz/backend-trainee-assignment/internal/app/validator"
	"github.com/psxzz/backend-trainee-assignment/internal/config"
)

type App struct {
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
	app.svc = service.New(storage, app.cfg.LogsPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't create a service: %w", err)
	}

	app.endp = endpoint.New(app.svc)

	app.echo = echo.New()
	app.echo.Validator = validator.New()

	// TODO: Declare endpoint handlers here
	app.echo.POST("/create", app.endp.HandleCreate)
	app.echo.DELETE("/delete", app.endp.HandleDelete)
	app.echo.POST("/experiments", app.endp.HandleExperiments)
	app.echo.GET("/list", app.endp.HandleUserExperimentList)
	app.echo.GET("/log/create", app.endp.HandleCreateLog)
	app.echo.GET("/log/download", app.endp.HandleDownloadLog)

	return app, nil
}

func (a App) Run() {
	go func() {
		if err := a.echo.Start(":8080"); err != nil && err != http.ErrServerClosed {
			a.echo.Logger.Fatal(err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //nolint:gomnd
	defer cancel()

	err := a.echo.Shutdown(ctx)
	if err != nil {
		a.echo.Logger.Fatal(err)
	}
}
