package repository

import (
	"backend/internal"
	"context"
	"os"
)

const (
	problemsDir = "../../problems/"
)

type Repo interface {
	GetFeaturedProblems(ctx context.Context) ([]internal.Problem, error)
	GetProblemText(ctx context.Context, id string) (string, error)
}

type repo struct {
	log internal.Logger
}

func NewRepo(logger internal.Logger) Repo {
	return &repo{
		log: logger,
	}
}

func (r repo) GetFeaturedProblems(_ context.Context) ([]internal.Problem, error) {
	entries, err := os.ReadDir(problemsDir)
	if err != nil {
		return nil, err
	}

	problems := make([]internal.Problem, 0, len(entries))
	for _, e := range entries {
		problems = append(problems, internal.Problem{
			Id:       e.Name(),
			Filename: e.Name(),
		})
	}
	return problems, nil
}

func (r repo) GetProblemText(_ context.Context, id string) (string, error) {
	content, err := os.ReadFile(problemsDir + id)
	return string(content), err
}
