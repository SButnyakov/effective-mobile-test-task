package test_helpers

import (
	"database/sql"
	"effective-mobile-test-task/internal/models"
	"math/rand"
	"time"
)

type generator struct {
	seed *rand.Rand
}

func GeneratePerson(delta int64) models.Person {
	g := generator{seed: rand.New(rand.NewSource(time.Now().UnixNano() + delta))}
	return models.Person{
		Name:        g.generateString(7),
		Surname:     g.generateString(7),
		Patronymic:  sql.NullString{String: g.generateString(7), Valid: true},
		Age:         g.generateAge(),
		Gender:      g.generateUint8Gender(),
		Nationality: g.generateString(2),
	}
}

func GeneratePersons(size int) []models.Person {
	persons := make([]models.Person, size)
	for i, _ := range persons {
		person := GeneratePerson(int64(i))
		person.Id = i + 1
		persons[i] = person
	}
	return persons
}

func (g *generator) generateString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = letters[g.seed.Intn(len(letters))]
	}
	return string(bytes)
}

func (g *generator) generateAge() uint8 {
	return uint8(g.seed.Intn(100))
}

func (g *generator) generateUint8Gender() uint8 {
	return uint8(g.seed.Intn(2))
}
