package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"` // ObjectID otomatis dari MongoDB
	Name        string             `bson:"name" json:"name"`
	Latitude    float64            `bson:"latitude" json:"latitude"`
	Longitude   float64            `bson:"longitude" json:"longitude"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
}