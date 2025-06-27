package repository

import (
	"context"
	"shop-backend/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, updated *model.Product) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context) ([]*model.Product, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Product, error)
}

type productRepo struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) ProductRepository {
	return &productRepo{
		collection: db.Collection("products"),
	}
}

func (r *productRepo) Create(ctx context.Context, product *model.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *productRepo) Update(ctx context.Context, updated *model.Product) error {
	filter := bson.M{"_id": updated.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        updated.Name,
			"description": updated.Description,
			"price":       updated.Price,
			"quantity":    updated.Stock,
			"image_url":   updated.ImageURL,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *productRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *productRepo) List(ctx context.Context) ([]*model.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*model.Product
	for cursor.Next(ctx) {
		var p model.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

func (r *productRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Product, error) {
	var product model.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
