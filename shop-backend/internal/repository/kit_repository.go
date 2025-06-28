package repository

import (
	"context"
	"shop-backend/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type KitRepository interface {
	Create(ctx context.Context, kit *model.Kit) error
}

type kitRepo struct {
	collection *mongo.Collection
}

func NewKitRepository(db *mongo.Database) KitRepository {
	return &kitRepo{
		collection: db.Collection("kits"),
	}
}

func (r *kitRepo) Create(ctx context.Context, kit *model.Kit) error {
	_, err := r.collection.InsertOne(ctx, kit)
	return err
}
