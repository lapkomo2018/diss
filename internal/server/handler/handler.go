package handler

import (
	"net/http"

	v1 "lapkomo2018/diss/internal/server/handler/v1"

	"github.com/go-chi/chi/v5"
)

type Params struct {
	V1 v1.Params
}

func Init(p Params) (http.Handler, error) {
	r := chi.NewRouter()

	v1Handler, err := v1.Init(p.V1)
	if err != nil {
		return nil, err
	}
	r.Mount("/v1", v1Handler)

	return r, nil
}
