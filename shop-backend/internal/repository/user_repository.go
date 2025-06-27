package repository

import (
	"context"
	"errors"
	"log"
	"shop-backend/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepo{
		collection: db.Collection("users"),
	}
}

var ErrUserNotFound = errors.New("user not found")

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Email not found, return nil user and nil error
			return nil, nil
		}
		// Some other error occurred
		return nil, err
	}
	return &user, nil
}

// Update user (used for OTP verification and general updates)
func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	filter := bson.M{"email": user.Email}
	updateFields := bson.M{
		"otp":         user.Otp,
		"otp_expiry":  user.OtpExpiry,
		"is_verified": user.IsVerified,
		"updated_at":  user.UpdatedAt,
	}
	// Optionally update password and phone if provided
	if user.Password != "" {
		updateFields["password"] = user.Password
	}
	if user.Phone != "" {
		updateFields["phone"] = user.Phone
	}

	update := bson.M{"$set": updateFields}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("no user found to update")
	}

	log.Printf("Updating user %s with fields: %+v", user.Email, updateFields)

	return nil
}
