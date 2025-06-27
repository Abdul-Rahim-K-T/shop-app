package database

import (
	"context"
	"log"
	"time"

	"shop-backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Database {
	cfg := config.LoadConfig()

	log.Printf("Attempting to connect to MongoDB at: %s", cfg.MongoURI) // Debuggging

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	log.Println("✅ Connected to MongoDB")

	// Replace 'shopdb' with you actual DB name
	return client.Database(cfg.DBName)
}
