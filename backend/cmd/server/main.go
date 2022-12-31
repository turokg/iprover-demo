package main

import (
	"backend/internal/conf"
	"backend/internal/ws"
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

func main() {
	flag.Parse()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(w, r)
	})
	err := http.ListenAndServe(conf.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type Args struct {
	inputParams string // passed as stdin to launched programs
}

func parseArgs(r *http.Request) (*Args, error) {
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
	return &Args{
		inputParams: params["message"][0],
	}, err
}

func checkOrigin(r *http.Request) bool {
	// TODO check origin properly
	return true
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
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
	client := ws.NewClient(ctx, conn)
	args, err := parseArgs(r)
	if err != nil {
		log.Println("something wrong", err)
		return
	}
	go client.Start(args.inputParams)
}
