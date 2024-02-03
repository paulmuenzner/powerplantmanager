package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserAuth struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Email           string             `bson:"email" json:"email" validate:"required,email,max=255" unique:"true"`
	Password        string             `bson:"password" json:"password,omitempty"`
	Verified        bool               `bson:"verified" json:"verified" default:"false"`
	VerifyToken     string             `bson:"verify_token" json:"verify_token"`
	Blocked         bool               `bson:"blocked" json:"blocked" default:"false"`
	DateVerifyToken time.Time          `bson:"date_verify_token" json:"date_verify_token"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at" validate:"required"`
}
