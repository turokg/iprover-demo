package launch

import (
	"backend/internal"
	"backend/internal/api"
	"backend/internal/launcher"
	"backend/internal/repository"
	"backend/internal/ws"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
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
	if len(params["problemId"]) != 1 {
		return nil, fmt.Errorf("couldn't parse query params")
	}
	return &launcher.LaunchArgs{
		ProblemID: params["problemId"][0],
	}, err
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

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = checkOrigin
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		msg := "unable to upgrade to webSockets"
		w.WriteHeader(internal.StatusError)
		w.Write([]byte(msg))
		h.logger.Error(ctx, msg, err)
		return
	}

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

	l := launcher.NewLauncher(h.logger)
	msgs, err := l.Launch(ctx, *args)
	if err != nil {
		msg := "unable to launch the process"
		w.WriteHeader(internal.StatusError)
		w.Write([]byte(msg))
		h.logger.Error(ctx, msg, err)
		return
	}
	client := ws.NewClient(ctx, conn, msgs, h.logger, cancel)
	go client.Start()
}
