package api

import (
	"effective-mobile-test-task/internal/http/api/v1"
	"effective-mobile-test-task/internal/http/api/v1/persons"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func Versions(log *slog.Logger, service persons.PersonService) chi.Router {
	router := chi.NewRouter()
	router.Mount("/v1", v1.Router(log, service))
	return router
}
