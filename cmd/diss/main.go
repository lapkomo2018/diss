package main

import (
	"log"
	"time"

	"lapkomo2018/diss/internal/server"
	"lapkomo2018/diss/internal/server/handler"
	v1 "lapkomo2018/diss/internal/server/handler/v1"
	"lapkomo2018/diss/internal/service/chat"
)

func main() {
	s := server.NewServer(":8080")
	if err := s.Init(handler.Params{
		V1: v1.Params{
			Chat: chat.NewChat(time.Millisecond*100, 10),
		},
	}); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	log.Printf("Starting server on %s", s.Server().Addr)
	if err := s.Start(); err != nil {
		panic(err)
	}
}
