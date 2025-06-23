package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ariefsibuea/articles-feed/internal/api/handler"
	"github.com/ariefsibuea/articles-feed/internal/api/repository"
	"github.com/ariefsibuea/articles-feed/internal/api/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	cfg := getConfig()
	e := echo.New()

	// set the minimum level of log
	e.Logger.SetLevel(log.INFO)

	// server timeout configurations
	e.Server.ReadTimeout = cfg.HTTPReadTimeout
	e.Server.WriteTimeout = cfg.HTTPWriteTimeout
	e.Server.IdleTimeout = cfg.HTTPIdleTimeout

	// customize error handler
	e.HTTPErrorHandler = handler.ErrorHandler()

	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		e.Logger.Fatalf("unable to parse database config: %v", err)
	}

	poolConfig.MaxConns = cfg.DBMaxConns
	poolConfig.MaxConnLifetime = cfg.DBMaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.DBMaxConnIdleTime
	poolConfig.MinConns = cfg.DBMinConns
	poolConfig.HealthCheckPeriod = cfg.DBHealthcheckPeriod

	dbpool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		e.Logger.Fatalf("unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		e.Logger.Fatalf("unable to ping database: %v", err)
	}

	// healthcheck endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "healthy",
		})
	})

	// init repositories
	articleRepository := repository.InitArticleRepository(dbpool)
	authorRepository := repository.InitAuthorRepository(dbpool)

	// init usecase
	articleUseCase := usecase.InitArticleUseCase(articleRepository, authorRepository)

	// init handler
	handler.InitArticleHandler(e, articleUseCase)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server")
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
