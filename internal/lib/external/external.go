package external

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type APIs struct {
}

func (a APIs) get(url string, v any) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if err = render.DecodeJSON(res.Body, v); err != nil {
		return err
	}

	if err = validator.New().Struct(v); err != nil {
		return err
	}

	return nil
}
