package persons

import (
	"effective-mobile-test-task/internal/lib/logger"
	"effective-mobile-test-task/internal/repositories"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"strconv"
)

// @Summary Delete
// @Tags persons
// @Description delete person
// @ID delete-person
// @Param id query int true "person id"
// @Router /persons/delete [post]
func Delete(log *slog.Logger, service PersonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "deleteOne"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Error("invalid id", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = service.Delete(id)
		if errors.Is(repositories.ErrPersonNotFound, err) {
			log.Info("person not found", slog.Int("id", id))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			log.Error("failed to delete person", logger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
