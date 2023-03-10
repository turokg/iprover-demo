package api

import (
	"backend/internal"
	"backend/internal/launcher"
	"backend/internal/ws"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

func parseArgs(r *http.Request) (*launcher.LaunchArgs, error) {
	ur, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		return nil, err
	}
	params, err := url.ParseQuery(ur.RawQuery)
	if err != nil {
		return nil, err
	}
	if len(params["message"]) != 1 {
		return nil, fmt.Errorf("couldn't parse query params")
	}
	return &launcher.LaunchArgs{
		Filename: params["message"][0],
	}, err
}

func checkOrigin(r *http.Request) bool {
	// TODO check origin properly
	return true
}

func NewWsHandler(logger internal.Logger) Handler {
	return &wsHandler{
		logger: logger,
	}
}

type wsHandler struct {
	logger internal.Logger
}

func (h *wsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = checkOrigin
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()

	args, err := parseArgs(r)
	if err != nil {
		log.Println("something wrong", err)
		return
	}

	l := launcher.NewLauncher(h.logger)
	msgs, err := l.Launch(ctx, *args)
	if err != nil {
		log.Println("couldn't launch", err)
		return
	}
	client := ws.NewClient(ctx, conn, msgs, h.logger)
	go client.Start()
}
