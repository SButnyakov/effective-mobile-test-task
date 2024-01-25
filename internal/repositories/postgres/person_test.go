package postgres

import (
	"database/sql"
	thelpers "effective-mobile-test-task/internal/lib/test_helpers"
	"effective-mobile-test-task/internal/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (p *PostgresTestSuite) TestPersonRepository_Create() {
	repo := NewPersonRepository(p.db)

	noNamePerson := thelpers.GeneratePerson(int64(1))
	noNamePerson.Name = ""

	noSurnamePerson := thelpers.GeneratePerson(int64(2))
	noSurnamePerson.Surname = ""

	noPatronymicPerson := thelpers.GeneratePerson(int64(3))
	noPatronymicPerson.Patronymic = sql.NullString{String: "", Valid: false}

	noNationalityPerson := thelpers.GeneratePerson(int64(4))
	noNationalityPerson.Nationality = ""

	tests := []struct {
		name  string
		input models.Person
		err   bool
	}{
		{"no name", noNamePerson, true},
		{"no surname", noSurnamePerson, true},
		{"no patronymic", noPatronymicPerson, false},
		{"no nationality", noNationalityPerson, true},
		{"all good", thelpers.GeneratePerson(int64(5)), false},
	}

	for _, t := range tests {
		if t.err {
			assert.Error(p.T(), repo.Create(t.input), t.name)
		} else {
			assert.NoError(p.T(), repo.Create(t.input), t.name)
		}
	}
}

func (p *PostgresTestSuite) TestPersonRepository_GetOne() {
	repo := NewPersonRepository(p.db)

	personWithPatronymic := thelpers.GeneratePerson(int64(1))
	personWithPatronymic.Id = 1

	personNoPatronymic := thelpers.GeneratePerson(int64(2))
	personNoPatronymic.Id = 2
	personNoPatronymic.Patronymic = sql.NullString{}

	persons := []models.Person{
		personWithPatronymic,
		personNoPatronymic,
	}

	p.seedPersons(repo, persons)

	resPatronymic, err := repo.GetOne(1)
	assert.NoError(p.T(), err)
	assert.Equal(p.T(), personWithPatronymic, resPatronymic)

	resNoPatronymic, err := repo.GetOne(2)
	assert.NoError(p.T(), err)
	assert.Equal(p.T(), personNoPatronymic, resNoPatronymic)
}

func (p *PostgresTestSuite) TestPersonRepository_GetMany() {
	repo := NewPersonRepository(p.db)

	persons := thelpers.GeneratePersons(20)

	p.seedPersons(repo, persons)

	res, err := repo.GetMany("SELECT * FROM persons")
	require.NoError(p.T(), err)

	for i, v := range res {
		fmt.Println(persons[i], v)
		assert.Equal(p.T(), persons[i], v)
	}
}

func (p *PostgresTestSuite) TestPersonRepository_Delete() {
	repo := NewPersonRepository(p.db)

	person := thelpers.GeneratePerson(0)
	person.Id = 1

	p.seedPersons(repo, []models.Person{person})

	require.NoError(p.T(), repo.Delete(1))
	_, err := repo.GetOne(1)
	require.Error(p.T(), err)
}

func (p *PostgresTestSuite) TestPersonRepository_Update() {
	repo := NewPersonRepository(p.db)

	person := thelpers.GeneratePerson(int64(1))
	person.Id = 1

	p.seedPersons(repo, []models.Person{person})

	res, err := repo.GetOne(person.Id)
	require.NoError(p.T(), err)
	require.Equal(p.T(), person, res)

	assert.Error(p.T(), repo.Update("UPDATE persons SET name=$1 WHERE id=$2", "", person.Id))
	assert.Error(p.T(), repo.Update("UPDATE persons SET surname=$1 WHERE id=$2", "", person.Id))
	assert.NoError(p.T(), repo.Update("UPDATE persons SET patronymic=$1 WHERE id=$2", sql.NullString{}, person.Id))
	assert.Error(p.T(), repo.Update("UPDATE persons SET nationality=$1 WHERE id=$2", "", person.Id))

	person = models.Person{
		Id:          person.Id,
		Name:        "name",
		Surname:     "surname",
		Patronymic:  sql.NullString{},
		Age:         27,
		Gender:      1,
		Nationality: "AA",
	}

	query := "UPDATE persons SET name=$1, surname=$2, patronymic=$3, age=$4, gender=$5, nationality=$6 WHERE id=$7"
	require.NoError(p.T(), repo.Update(query, person.Name, person.Surname, person.Patronymic, person.Age,
		person.Gender, person.Nationality, person.Id))

	res, err = repo.GetOne(person.Id)
	require.NoError(p.T(), err)
	require.Equal(p.T(), person, res)
}

func (p *PostgresTestSuite) seedPersons(repo *PersonRepository, persons []models.Person) {
	for _, person := range persons {
		require.NoError(p.T(), repo.Create(person))
	}
}
