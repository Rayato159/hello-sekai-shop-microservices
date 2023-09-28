package player

import "time"

type (
	PlayerProfile struct {
		Id        string    `json:"_id"`
		Email     string    `json:"email"`
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
