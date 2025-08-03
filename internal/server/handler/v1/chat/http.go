package chat

import (
	"io"
	"net/http"

	"lapkomo2018/diss/internal/service/chat"
)

func publishHandler(chat chat.Chat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := http.MaxBytesReader(w, r.Body, 8192)
		msg, err := io.ReadAll(body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
			return
		}

		if err := chat.Publish(r.Context(), msg); err != nil {
			http.Error(w, "failed to publish message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
