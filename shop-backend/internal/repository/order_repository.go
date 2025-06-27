package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository struct {
	DB *mongo.Database
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{DB: db}
}
