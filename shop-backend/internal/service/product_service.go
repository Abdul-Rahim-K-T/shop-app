package service

import (
	"context"
	"shop-backend/internal/model"
	"shop-backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	Repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *model.Product) error {
	return s.Repo.Create(ctx, product)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *model.Product) error {
	return s.Repo.Update(ctx, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	return s.Repo.Delete(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context) ([]*model.Product, error) {
	return s.Repo.List(ctx)
}

func (s *ProductService) GetByIDProduct(ctx context.Context, id primitive.ObjectID) (*model.Product, error) {
	return s.Repo.FindByID(ctx, id)
}
