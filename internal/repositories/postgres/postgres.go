package postgres

import (
	"database/sql"
	"effective-mobile-test-task/internal/lib/config"
	"fmt"
	_ "github.com/lib/pq"
)

func New(cfg config.PG) (*sql.DB, error) {
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return db, nil
}
