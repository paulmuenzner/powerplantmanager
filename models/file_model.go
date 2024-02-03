package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

type File struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	PublicFileID string             `bson:"public_file_id" json:"public_file_id" validate:"required" unique:"true"`
	Name         string             `bson:"name" json:"name" validate:"required,max=130" unique:"false"`
	Slug         string             `bson:"slug" json:"slug" validate:"required,max=130" unique:"true"`
	Type         string             `bson:"type" json:"type" validate:"required"` // file type (.jpeg, .png, .pdf)
	User         primitive.ObjectID `bson:"user_id"`                              // _id of user who owns/uploaded this file
	Size         int                `bson:"size" json:"size"  validate:"required" unique:"false"`
	Width        int                `bson:"width" json:"width"  validate:"required" unique:"false"`
	Height       int                `bson:"height" json:"height"  validate:"required" unique:"false"`
	Folder       string             `bson:"folder" validate:"required" unique:"false"` // folder in bucket
	CreatedAt    time.Time          `bson:"created_at" json:"created_at" validate:"required"`
}

func FileType(fl validator.FieldLevel) bool {
	allowedTypes := map[string]bool{
		"jpeg": true,
		"png":  true,
		"pdf":  true,
		// Add more allowed types as needed
	}

	fileType := strings.ToLower(fl.Field().String())
	return allowedTypes[fileType]
}
