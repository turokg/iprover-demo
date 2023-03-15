package launch

import (
	"backend/internal"
	"backend/internal/api"
	"backend/internal/launcher"
	"backend/internal/repository"
	"backend/internal/ws"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

func parseArgs(r *http.Request) (launcher.LaunchArgs, error) {
	args := launcher.LaunchArgs{}
	ur, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		return args, err
	}
	params, err := url.ParseQuery(ur.RawQuery)
	if err != nil {
		return args, err
	}
	if len(params["problemId"]) != 1 {
		return args, fmt.Errorf("couldn't parse query params")
	}
	args.ProblemID = params["problemId"][0]
	return args, err
}

func checkOrigin(r *http.Request) bool {
	// TODO check origin properly
	return true
}

type handler struct {
	logger internal.Logger
	repo   repository.Repo
}

func New(logger internal.Logger, repo repository.Repo) api.Handler {
	return &handler{
		logger: logger,
		repo:   repo,
	}
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), internal.RunTimeout)

	args, err := parseArgs(r)
	if err != nil {
		msg := "unable to parse args"
		w.WriteHeader(internal.StatusError)
		w.Write([]byte(msg))
		h.logger.Error(ctx, msg, err)
		return
	}

	args.ProblemText, err = h.repo.GetProblemText(ctx, args.ProblemID)
	if err != nil {
		msg := "unable to fetch problem text"
		w.WriteHeader(internal.StatusError)
		w.Write([]byte(msg))
		h.logger.Error(ctx, msg, err)
		return
	}
	output := make(chan []byte, internal.LaunchBuffer)
	client, err := ws.NewClient(w, r, output, h.logger)
	if err != nil {
		msg := "unable to upgrade to webSockets"
		w.WriteHeader(internal.StatusError)
		w.Write([]byte(msg))
		h.logger.Error(ctx, msg, err)
		return
	}

	l := launcher.NewLauncher(h.logger)

	wg := &sync.WaitGroup{}
	go l.Launch(ctx, wg, args, output)
	go client.Read(ctx, wg, cancel)
	go client.Write(ctx, wg)
	wg.Wait()
	// TODO почему не ждет то???

	h.logger.Info(ctx, "handler finished")
}
