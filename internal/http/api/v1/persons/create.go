package persons

import (
	"effective-mobile-test-task/internal/lib/converters"
	"effective-mobile-test-task/internal/lib/logger"
	"effective-mobile-test-task/internal/models"
	"effective-mobile-test-task/internal/services"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type CreateOneRequest struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic,omitempty"`
}

// @Summary Create
// @Tags persons
// @Description create person
// @ID create-person
// @Accept json
// @Param input body CreateOneRequest true "person's info with patronymic"
// @Router /persons/create [post]
func Create(log *slog.Logger, service PersonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "createOne"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req CreateOneRequest

		err := render.DecodeJSON(r.Body, &req)
		log.Debug("request", slog.Any("req", req))
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = validator.New().Struct(req); err != nil {
			log.Error("invalid request", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Debug("decoded request:", slog.Any("person", req))

		err = service.Create(models.Person{
			Name:       req.Name,
			Surname:    req.Surname,
			Patronymic: converters.StringToNullString(req.Patronymic),
		})
		if errors.Is(err, services.ErrExternalError) {
			// logged
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		if err != nil {
			log.Error("failed to store new person", logger.Err(err))
			w.WriteHeader(http.StatusInsufficientStorage)
		}

		w.WriteHeader(http.StatusCreated)
	}
}
