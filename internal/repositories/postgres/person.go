package postgres

import (
	"database/sql"
	"effective-mobile-test-task/internal/models"
	"effective-mobile-test-task/internal/repositories"
	"errors"
	"fmt"
)

const packagePath = "storage.postgres.repositories."

type PersonRepository struct {
	db *sql.DB
}

func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{
		db: db,
	}
}

func (p *PersonRepository) Create(person models.Person) error {
	const fn = packagePath + "Create"

	query := "INSERT INTO persons (name, surname, patronymic, age, gender, nationality) " +
		"VALUES ($1, $2, $3, $4, $5, $6)"

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}

func (p *PersonRepository) GetOne(id int) (models.Person, error) {
	const fn = packagePath + "GetOne"

	query := "SELECT id, name, surname, patronymic, age, gender, nationality FROM persons WHERE id = $1"

	var person models.Person

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return person, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	err = stmt.QueryRow(id).Scan(&person.Id, &person.Name, &person.Surname, &person.Patronymic,
		&person.Age, &person.Gender, &person.Nationality)
	if errors.Is(err, sql.ErrNoRows) {
		return person, repositories.ErrPersonNotFound
	}
	if err != nil {
		return person, fmt.Errorf("%s: scanning row: %w", fn, err)
	}

	return person, nil
}

func (p *PersonRepository) GetMany(query string, a ...any) ([]models.Person, error) {
	const fn = packagePath + "GetMany"

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	rows, err := stmt.Query(a...)
	if err != nil {
		return nil, fmt.Errorf("%s: executing query: %w", fn, err)
	}
	defer rows.Close()

	persons := make([]models.Person, 0)

	for rows.Next() {
		var person models.Person
		err = rows.Scan(&person.Id, &person.Name, &person.Surname, &person.Patronymic,
			&person.Age, &person.Gender, &person.Nationality)
		if err != nil {
			return nil, fmt.Errorf("%s: scanning rows: %w", fn, err)
		}
		persons = append(persons, person)
	}

	return persons, nil
}

func (p *PersonRepository) Delete(id int) error {
	const fn = packagePath + "Delete"

	query := "DELETE FROM persons WHERE id=$1"

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: get affected rows: %w", fn, err)
	} else if rows == 0 {
		return repositories.ErrPersonNotFound
	}

	return nil
}

func (p *PersonRepository) Update(query string, a ...any) error {
	const fn = packagePath + "Update"

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	res, err := stmt.Exec(a...)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: get affected rows: %w", fn, err)
	} else if rows == 0 {
		return repositories.ErrPersonNotFound
	}

	return nil
}
