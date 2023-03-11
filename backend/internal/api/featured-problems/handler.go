package featured_problems

import (
	"backend/internal"
	"backend/internal/api"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"net/http"
)

func New(logger internal.Logger, repo repository.Repo) api.Handler {
	return &handler{
		logger: logger,
		repo:   repo,
	}
}

type handler struct {
	logger internal.Logger
	repo   repository.Repo
}

type dto struct {
	Id       string `json:"id,omitempty"`
	Filename string `json:"filename,omitempty"`
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	problems, err := h.repo.GetFeaturedProblems(ctx)
	if err != nil {
		w.WriteHeader(internal.StatusError) // TODO change to constant
		w.Write([]byte("unable to get the list of problems"))
		h.logger.Error(ctx, "unable to get the list of problems", err)
	}
	dtos := make([]dto, 0, len(problems))
	for _, p := range problems {
		dtos = append(dtos, dto{
			Id:       p.ID,
			Filename: p.Filename,
		})
	}
	resp, err := json.Marshal(dtos)
	if err != nil {
		h.logger.Error(ctx, "unable to marshall response", err)
	}

	_, err = w.Write(resp)
	if err != nil {
		h.logger.Error(ctx, "unable to write response", err)
	}
	return
}
