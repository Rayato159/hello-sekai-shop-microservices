package inventory

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	Inventory struct {
		Id       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		PlayerId string             `json:"player_id" bson:"player_id"`
		ItemId   string             `json:"item_id" bson:"item_id"`
	}
)
