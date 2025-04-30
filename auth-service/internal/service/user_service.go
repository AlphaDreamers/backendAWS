package service

import "github.com/SwanHtetAungPhyo/auth/internal/repo"

type Service interface {
}

type UserService struct {
	repo repo.Repository
}

func NewUserService(repo repo.Repository) Service {
	return &UserService{
		repo: repo,
	}
}
