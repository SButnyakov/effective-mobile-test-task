package v1

import (
	"effective-mobile-test-task/internal/http/api/v1/persons"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

// @title Effective Mobile API
// @version 1.0
// @description API HTTPServer for Effective Mobile Test Task

// @host localhost:8080
// @BasePath /api/v1

func Router(log *slog.Logger, service persons.PersonService) chi.Router {
	router := chi.NewRouter()
	router.Mount("/persons", persons.Router(log, service))
	return router
}
