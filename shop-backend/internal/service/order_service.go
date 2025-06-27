package service

import (
	"context"
	"shop-backend/internal/repository"
)

type OrderService struct {
	OrderRepo *repository.OrderRepository
}

func NewOrderService(orderRepo *repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepo: orderRepo,
	}
}

// Example method
func (s *OrderService) PlaceOrder(ctx context.Context, userID int, items []string) error {
	// Business logic here
	return nil
}
