package service

import (
	"context"
	"shop-backend/internal/model"
	"shop-backend/internal/repository"
)

type KitService struct {
	Repo repository.KitRepository
}

func NewKitService(repo repository.KitRepository) *KitService {
	return &KitService{Repo: repo}
}

func (s *KitService) CreateKit(ctx context.Context, kit *model.Kit) error {
	return s.Repo.Create(ctx, kit)
}
