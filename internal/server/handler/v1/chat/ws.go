package chat

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"lapkomo2018/diss/internal/service/chat"

	"github.com/coder/websocket"
)

func websocketHandler(chat chat.Chat) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			http.Error(w, "failed to accept websocket connection", http.StatusInternalServerError)
			return
		}
		defer conn.CloseNow()

		handleWebSocketConnection(r.Context(), chat, conn)
	}
}

func handleWebSocketConnection(ctx context.Context, chat chat.Chat, conn *websocket.Conn) {
	var (
		mu     sync.Mutex
		closed bool
	)

	sub := chat.NewSubscriber(func() {
		mu.Lock()
		defer mu.Unlock()
		if closed {
			return
		}
		closed = true
		conn.Close(websocket.StatusNormalClosure, "subscriber closed")
	})
	defer chat.Unsubscribe(sub)

	mu.Lock()
	if closed {
		mu.Unlock()
		log.Printf("websocket connection already closed")
		return
	}
	mu.Unlock()

	listenToMessages(ctx, sub, conn)
}

func listenToMessages(ctx context.Context, sub chat.Subscriber, conn *websocket.Conn) {
	ctx = conn.CloseRead(ctx)

	for {
		select {
		case msg, ok := <-sub.Messages():
			if !ok {
				log.Printf("subscriber messages channel closed")
				return
			}

			if err := writeTimeout(ctx, time.Second*5, conn, msg); err != nil {
				log.Printf("write error: %v", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
