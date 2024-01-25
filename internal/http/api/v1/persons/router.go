package persons

import (
	"effective-mobile-test-task/internal/models"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

const packagePath = "http.api.v1.persons."

type PersonService interface {
	Create(person models.Person) error
	Delete(id int) error
	Update(args map[string]interface{}) error
	GetOne(id int) (models.Person, error)
	GetMany(args map[string]interface{}) ([]models.Person, error)
}

func Router(log *slog.Logger, service PersonService) chi.Router {
	router := chi.NewRouter()
	router.Get("/", Get(log, service))
	router.Post("/create", Create(log, service))
	router.Post("/delete", Delete(log, service))
	router.Put("/update", Update(log, service))
	return router
}
