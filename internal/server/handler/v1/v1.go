package v1

import (
	"net/http"

	chatHandler "lapkomo2018/diss/internal/server/handler/v1/chat"
	"lapkomo2018/diss/internal/service/chat"

	"github.com/go-chi/chi/v5"
)

type Params struct {
	Chat chat.Chat
}

func Init(p Params) (http.Handler, error) {
	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	handler, err := chatHandler.Handler(p.Chat)
	if err != nil {
		return nil, err
	}
	r.Mount("/chat", handler)

	return r, nil
}
