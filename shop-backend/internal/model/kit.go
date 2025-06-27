package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Kit struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	ProductIDs  []primitive.ObjectID `bson:"product_ids" json:"product_ids"`
	Price       float64              `bson:"price" json:"price"`
	ImageURL    string               `bson:"image_url" json:"image_url"`
}
