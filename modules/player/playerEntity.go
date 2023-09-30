package player

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Player struct {
		Id          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Email       string             `json:"email" bson:"email"`
		Password    string             `json:"password" bson:"password"`
		Username    string             `json:"username" bson:"username"`
		CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
		UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
		PlayerRoles []PlayerRole       `bson:"player_roles"`
	}

	PlayerRole struct {
		RoleTitle string `json:"role_title" bson:"role_title"`
		RoleCode  int    `json:"role_code" bson:"role_code"`
	}

	PlayerProfileBson struct {
		Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Email     string             `json:"email" bson:"email"`
		Username  string             `json:"username" bson:"username"`
		CreatedAt time.Time          `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	}

	PlayerSavingAccount struct {
		PlayerId string  `json:"player_id" bson:"player_id"`
		Balance  float64 `json:"balance" bson:"balance"`
	}

	PlayerTransaction struct {
		Id        primitive.ObjectID `bson:"_id,omitempty"`
		PlayerId  string             `bson:"player_id"`
		Amount    float64            `bson:"amount"`
		CreatedAt time.Time          `bson:"created_at"`
	}
)
