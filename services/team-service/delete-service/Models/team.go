package team_delete_model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Team struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name       string             `json:"name" validate:"required" bson:"name"`
	Captain    string             `json:"captain" validate:"required" bson:"captain"`
	Players    []string           `json:"players" bson:"players"`
	HomeGround string             `json:"homeground" bson:"homeground"`
	Points     int                `json:"points,omitempty" bson:"points,truncate"`
}
