package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Base struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Enabled   bool               `json:"enabled" bson:"enabled"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt time.Time          `json:"deletedAt" bson:"deletedAt"`
}

func NewBase() Base {
	return Base{
		ID:        newID(),
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func newID() primitive.ObjectID {
	return primitive.NewObjectID()
}
