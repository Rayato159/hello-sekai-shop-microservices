package inventory

type (
	Inventory struct {
		Id       string `json:"_id" bson:"_id,omitempty"`
		PlayerId string `json:"player_id" bson:"player_id"`
		ItemId   string `json:"item_id" bson:"item_id"`
	}
)
