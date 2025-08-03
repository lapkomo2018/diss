package chat

import (
	"errors"
	"net/http"

	"lapkomo2018/diss/internal/service/chat"

	"github.com/go-chi/chi/v5"
)

var ErrNoChatService = errors.New("chat service is not available")

func Handler(chat chat.Chat) (http.Handler, error) {
	if chat == nil {
		return nil, ErrNoChatService
	}

	r := chi.NewRouter()

	r.Get("/ws", websocketHandler(chat))
	r.Post("/publish", publishHandler(chat))

	return r, nil
}
