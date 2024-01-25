package models

import "database/sql"

type Person struct {
	Id          int
	Name        string
	Surname     string
	Patronymic  sql.NullString
	Age         uint8
	Gender      uint8
	Nationality string
}
