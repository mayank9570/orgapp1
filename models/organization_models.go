package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Organization struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Location string             `json:"location,omitempty" validate:"required"`
	City     string             `json:"city,omitempty" validate:"required"`
	State    string             `json:"state,omitempty" validate:"required"`
	Email    string             `json:"email,omitempty" validate:"required"`
}
