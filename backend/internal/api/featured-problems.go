package api

import (
	"backend/internal"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"net/http"
)

func NewFeaturedProblemsHandler(logger internal.Logger, repo repository.Repo) Handler {
	return &fpHandler{
		logger: logger,
		repo:   repo,
	}
}

type fpHandler struct {
	logger internal.Logger
	repo   repository.Repo
}

type FpDTO struct {
	id       string
	filename string
}

func (h *fpHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	problems, err := h.repo.GetFeaturedProblems(ctx)
	if err != nil {
		w.WriteHeader(500) // TODO change to constant
	}
	dto := make([]FpDTO, 0, len(problems))
	for _, p := range problems {
		dto = append(dto, FpDTO{
			id:       p.Id,
			filename: p.Filename,
		})
	}
	resp, err := json.Marshal(dto)
	if err != nil {
		// TODO
	}
	_, err = w.Write(resp)
	if err != nil {
		// TODO
	}
	return
}
