package main

import (
	"context"
	_ "effective-mobile-test-task/docs"
	"effective-mobile-test-task/internal/http/api"
	mwLogger "effective-mobile-test-task/internal/http/middleware/logger"
	"effective-mobile-test-task/internal/lib/config"
	"effective-mobile-test-task/internal/lib/external"
	"effective-mobile-test-task/internal/lib/logger"
	"effective-mobile-test-task/internal/repositories/postgres"
	"effective-mobile-test-task/internal/services"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	log2 "log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	envPath := os.Getenv("ENV_PATH")

	err := godotenv.Load(envPath)
	if err != nil {
		log2.Fatalf("env file not found. Err: %s", err.Error())
	}

	cfg, err := config.Load()
	if err != nil {
		log2.Fatalf("failed to create config. Err: %s", err.Error())
	}

	// Logger setup
	log := logger.New(cfg.Env)
	log.Info("logger initialized")
	log.Debug("logger debug mode enabled")

	// Postgres
	pg, err := postgres.New(cfg.PG)
	if err != nil {
		log.Error("failed to initialize database.", logger.Err(err))
		os.Exit(1)
	}
	defer pg.Close()

	// Repositories
	personRepository := postgres.NewPersonRepository(pg)

	// Services
	personService := services.NewPersonService(personRepository, external.APIs{}, log, *cfg)

	// Router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(mwLogger.New(log))

	router.Mount("/swagger", httpSwagger.WrapHandler)
	router.Mount("/api", api.Versions(log, personService))
	log.Debug("router initialized")

	// Server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error.", logger.Err(err))
			os.Exit(1)
		}
	}()
	log.Info("server started at",
		slog.String("host", cfg.HTTPServer.Host),
		slog.Int("port", cfg.HTTPServer.Port))

	<-interrupt
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed.", logger.Err(err))
		os.Exit(1)
	}
}
