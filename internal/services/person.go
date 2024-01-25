package services

import (
	"effective-mobile-test-task/internal/lib/config"
	"effective-mobile-test-task/internal/lib/converters"
	"effective-mobile-test-task/internal/lib/logger"
	"effective-mobile-test-task/internal/models"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

type PersonProvider interface {
	Create(person models.Person) error
	Delete(id int) error
	GetOne(id int) (models.Person, error)
	GetMany(query string, a ...any) ([]models.Person, error)
	Update(query string, a ...any) error
}

type DataProvider interface {
	GetAge(name string) (uint8, error)
	GetMostProbableCountry(name string) (string, error)
	GetGender(name string) (string, error)
}

type PersonService struct {
	pProvider PersonProvider
	dProvider DataProvider
	log       *slog.Logger
	cfg       config.Config
}

func NewPersonService(pProvider PersonProvider, dProvider DataProvider, log *slog.Logger,
	cfg config.Config) *PersonService {
	return &PersonService{
		pProvider: pProvider,
		dProvider: dProvider,
		log:       log,
		cfg:       cfg,
	}
}

func (s *PersonService) Create(person models.Person) error {
	s.log.Debug("adding", slog.Any("person", person))

	ageChan := make(chan uint8)
	nationalityChan := make(chan string)
	genderChan := make(chan string)

	s.log.Debug("starting goroutines")
	go s.getAgeSync(ageChan, person.Name)
	go s.getNationalityAsync(nationalityChan, person.Name)
	go s.getGenderAsync(genderChan, person.Name)

	age, ok := <-ageChan
	if !ok {
		return ErrExternalError
	}
	s.log.Debug("got", slog.Any("age", age))

	nationality, ok := <-nationalityChan
	if !ok {
		return ErrExternalError
	}
	s.log.Debug("got", slog.String("nationality", nationality))

	gender, ok := <-genderChan
	if !ok {
		return ErrExternalError
	}
	s.log.Debug("got", slog.Any("gender", gender))

	person.Age = age
	person.Nationality = nationality
	uint8Gender, err := converters.GenderStringToUint8(gender)
	if err != nil {
		return ErrExternalError
	}
	person.Gender = uint8Gender

	s.log.Debug("final person struct", slog.Any("person", person))

	return s.pProvider.Create(person)
}

func (s *PersonService) Delete(id int) error {
	s.log.Debug("deleting", slog.Int("id", id))
	return s.pProvider.Delete(id)
}

func (s *PersonService) Update(argsMap map[string]interface{}) error {
	query, args := generateUpdateQueryAndArgs(argsMap)
	s.log.Debug("generated", slog.String("query", query))
	s.log.Debug("generated", slog.Any("args", args))
	return s.pProvider.Update(query, args...)
}

func (s *PersonService) GetOne(id int) (models.Person, error) {
	return s.pProvider.GetOne(id)
}

func (s *PersonService) GetMany(argsMap map[string]interface{}) ([]models.Person, error) {
	query, args := generateGetManyQueryAndArgs(argsMap, s.cfg.API.PersonPageLimit)
	s.log.Debug("generated", slog.String("query", query))
	s.log.Debug("generated", slog.Any("args", args))
	return s.pProvider.GetMany(query, args...)
}

func (s *PersonService) getAgeSync(out chan<- uint8, name string) {
	age, err := s.dProvider.GetAge(name)
	if err != nil {
		s.log.Error("failed to get age.", logger.Err(err))
		close(out)
		return
	}
	out <- age
}

func (s *PersonService) getNationalityAsync(out chan<- string, name string) {
	nationality, err := s.dProvider.GetMostProbableCountry(name)
	if err != nil {
		s.log.Error("failed to get nationality", logger.Err(err))
		close(out)
		return
	}
	out <- nationality
}

func (s *PersonService) getGenderAsync(out chan<- string, name string) {
	gender, err := s.dProvider.GetGender(name)
	if err != nil {
		s.log.Error("failed to get gender", logger.Err(err))
		close(out)
		return
	}
	out <- gender
}

func generateUpdateQueryAndArgs(argsMap map[string]interface{}) (string, []interface{}) {
	var sb strings.Builder
	sb.WriteString("UPDATE persons SET")
	argsCounter := 1
	args := []interface{}{}

	var id int

	for key, value := range argsMap {
		if key == "id" {
			id = value.(int)
			continue
		}
		if value != nil && !reflect.ValueOf(value).IsNil() {
			sb.WriteString(fmt.Sprintf(" %s=$%d,", key, argsCounter))
			args = append(args, value)
			argsCounter++
		}
	}

	query := sb.String()
	query = query[:len(query)-1]
	query = fmt.Sprintf("%s WHERE id=$%d", query, argsCounter)
	args = append(args, id)

	return query, args
}

func generateGetManyQueryAndArgs(argsMap map[string]interface{}, pageLength int) (string, []interface{}) {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM persons WHERE")
	argsCounter := 1
	args := []interface{}{}

	var paginated bool
	if val, ok := argsMap["page"]; ok {
		sb.WriteString(fmt.Sprintf(" id > %d  AND", val.(int)*pageLength-pageLength))
		paginated = true
		delete(argsMap, "page")
	}

	for key, value := range argsMap {
		if value != nil && !reflect.ValueOf(value).IsNil() {
			sb.WriteString(fmt.Sprintf(" %s=$%d  AND", key, argsCounter))
			args = append(args, value)
			argsCounter++
		}
	}

	query := sb.String()
	query = query[:len(query)-5]

	if paginated {
		query = fmt.Sprintf("%s ORDER BY id ASC LIMIT %d", query, pageLength)
	}

	return query, args
}
