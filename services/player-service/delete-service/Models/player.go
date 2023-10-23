package cru_model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name       string             `json:"name" validate:"required"`
	Jno        string             `json:"jno" validate:"required"`
	Age        string             `json:"age" validate:"required"`
	Role       [2]string          `json:"role" validate:"required"`
	BatAvg     float64            `json:"batavg,omitempty" bson:"batavg,truncate"`
	StrikeRate float64            `json:"strikerate,omitempty" bson:"strikerate,truncate"`
	Econ       float64            `json:"econ,omitempty" bson:"econ,truncate"`
	Wickets    int                `json:"wickets,omitempty" bson:"wickets,truncate"`
	Matches    int                `json:"matches,omitempty" bson:"matches,truncate"`
	Runs       int                `json:"runs,omitempty" bson:"runs,truncate"`
}
