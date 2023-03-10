package main

import (
	"backend/internal"
	"backend/internal/api"
	"context"
	"log"
	"net/http"
)

func main() {

	logger := internal.NewLogger()
	ctx := context.Background()

	wsHandler := api.NewWsHandler(logger)
	http.HandleFunc("/ws", wsHandler.Handle)

	logger.Info(ctx, "starting we server")
	err := http.ListenAndServe(internal.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
