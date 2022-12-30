package main

import (
	"backend/internal/launcher"
	"context"
	"fmt"
	"log"
	"time"
)

const inputParams = "hello worldcsscs"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	messages, err := launcher.Launch(ctx, inputParams)
	if err != nil {
		log.Fatalln("couldn't start", err)
	}

	go func() {
		time.Sleep(time.Second * 2)
		cancel()
	}()
	for msg := range messages {
		fmt.Println(msg)
	}
}
