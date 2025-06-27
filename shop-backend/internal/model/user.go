package model

import "time"

type User struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	FirstName  string    `bson:"first_name" json:"first_name"`
	LastName   string    `bson:"last_name" json:"last_name"`
	Email      string    `bson:"email" json:"email"`
	Phone      string    `bson:"phone" json:"phone"`
	Password   string    `bson:"password" json:"-"`
	Address    string    `bson:"address" json:"address"`
	IsVerified bool      `bson:"is_verified" json:"is_verified"`
	Otp        string    `bson:"otp" json:"-"`
	OtpExpiry  time.Time `json:"otp_expiry,omitempty" bson:"otp_expiry,omitempty"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}
