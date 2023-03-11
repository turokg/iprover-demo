package main

import (
	"backend/internal"
	featured_problems "backend/internal/api/featured-problems"
	"backend/internal/api/launch"
	"backend/internal/repository"
	"context"
	"log"
	"net/http"
)

func main() {

	logger := internal.NewLogger()
	ctx := context.Background()

	repo := repository.New(logger)

	launchHandler := launch.New(logger, repo)
	http.HandleFunc("/launch", launchHandler.Handle)

	fpHandler := featured_problems.New(logger, repo)
	http.HandleFunc("/featured-problems", fpHandler.Handle)

	logger.Info(ctx, "starting we server")
	err := http.ListenAndServe(internal.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
