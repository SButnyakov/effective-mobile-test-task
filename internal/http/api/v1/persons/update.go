package persons

import (
	"effective-mobile-test-task/internal/lib/converters"
	"effective-mobile-test-task/internal/lib/logger"
	"effective-mobile-test-task/internal/repositories"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type UpdateOneRequest struct {
	Id          int     `json:"id" validate:"required"`
	Name        *string `json:"name" validate:"omitempty,min=1"`
	Surname     *string `json:"surname" validate:"omitempty,min=1"`
	Patronymic  *string `json:"patronymic" validate:"omitempty,min=1"`
	Age         *uint8  `json:"age" validate:"omitempty,min=0"`
	Gender      *string `json:"gender" validate:"omitempty,oneof=male female"`
	Nationality *string `json:"nationality" validate:"omitempty,min=1"`
}

// @Summary Update
// @Tags persons
// @Description update person
// @ID update-person
// @Accept json
// @Param input body UpdateOneRequest true "person's update info"
// @Router /persons/update [put]
func Update(log *slog.Logger, service PersonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "updateOne"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req UpdateOneRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Debug("decoded request:", slog.Any("person", req))
		if err = validator.New().Struct(req); err != nil {
			log.Error("invalid request", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = service.Update(argsToMapUpdate(req))
		if errors.Is(err, repositories.ErrPersonNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			log.Error("failed to update person", logger.Err(err))
			w.WriteHeader(http.StatusInsufficientStorage)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func argsToMapUpdate(request UpdateOneRequest) map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = request.Id
	m["name"] = request.Name
	m["surname"] = request.Surname
	m["patronymic"] = request.Patronymic
	m["age"] = request.Age
	if request.Gender == nil {
		m["gender"] = nil
	} else {
		uint8Gender, err := converters.GenderStringToUint8(*request.Gender)
		if err != nil {
			m["gender"] = nil
		} else {
			m["gender"] = &uint8Gender
		}
	}
	m["nationality"] = request.Nationality
	return m
}
