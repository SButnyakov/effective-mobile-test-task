package persons

import (
	"effective-mobile-test-task/internal/lib/converters"
	"effective-mobile-test-task/internal/lib/logger"
	"effective-mobile-test-task/internal/repositories"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type GetManyResponse struct {
	Persons []GetOneResponse `json:"persons"`
}

type GetOneResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         uint8  `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

// @Summary Get
// @Tags persons
// @Description get persons
// @ID get-persons
// @Param id query int false "person id"
// @Param name query string false "person name"
// @Param surname query string false "person surname"
// @Param patronymic query string false "person patronymic"
// @Param age query int false "person age"
// @Param gender query string false "person gender"
// @Param nationality query string false "person nationality"
// @Param page query int false "page"
// @Success 200 {object} GetManyResponse
// @Success 200 {object} GetOneResponse
// @Router /persons [get]
func Get(log *slog.Logger, service PersonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "Get"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := r.URL.Query().Get("id")

		switch id {
		case "":
			getMany(log, service, w, r)
		default:
			getOne(log, id, service, w, r)
		}
	}
}

func getOne(log *slog.Logger, id string, service PersonService, w http.ResponseWriter, r *http.Request) {
	iid, err := strconv.Atoi(id)
	if err != nil {
		log.Error("invalid id", logger.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	person, err := service.GetOne(iid)
	if errors.Is(err, repositories.ErrPersonNotFound) {
		log.Error("person not found", logger.Err(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error("failed to get person", logger.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := GetOneResponse{
		Id:          person.Id,
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  converters.NullStringToString(person.Patronymic),
		Age:         person.Age,
		Gender:      converters.GenderUint8ToString(person.Gender),
		Nationality: person.Nationality,
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, response)
}

func getMany(log *slog.Logger, service PersonService, w http.ResponseWriter, r *http.Request) {
	persons, err := service.GetMany(argsToMapGetMany(r))
	if errors.Is(err, repositories.ErrPersonNotFound) {
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, GetManyResponse{Persons: []GetOneResponse{}})
		return
	}
	if err != nil {
		log.Error("failed to get persons", logger.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responsePersons := make([]GetOneResponse, len(persons))
	for i, v := range persons {
		responsePersons[i] = GetOneResponse{
			Id:          v.Id,
			Name:        v.Name,
			Surname:     v.Surname,
			Patronymic:  converters.NullStringToString(v.Patronymic),
			Age:         v.Age,
			Gender:      converters.GenderUint8ToString(v.Gender),
			Nationality: v.Nationality,
		}
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, GetManyResponse{Persons: responsePersons})
}

func argsToMapGetMany(r *http.Request) map[string]interface{} {
	m := make(map[string]interface{})
	m["name"] = converters.StringToStringPointer(r.URL.Query().Get("name"))
	m["surname"] = converters.StringToStringPointer(r.URL.Query().Get("surname"))
	m["patronymic"] = converters.StringToStringPointer(r.URL.Query().Get("patronymic"))
	age, err := strconv.ParseUint(r.URL.Query().Get("age"), 10, 8)
	if err != nil {
		m["age"] = nil
	} else {
		m["age"] = &age
	}
	gender, err := converters.GenderStringToUint8(r.URL.Query().Get("gender"))
	if err != nil {
		m["gender"] = nil
	} else {
		m["gender"] = &gender
	}
	m["nationality"] = converters.StringToStringPointer(r.URL.Query().Get("nationality"))
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err == nil {
		m["page"] = page
	}
	return m
}
