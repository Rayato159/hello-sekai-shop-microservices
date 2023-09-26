package itemRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ItemRepositoryService interface{}

	itemRepository struct {
		db *mongo.Client
	}
)

func NewItemRepository(db *mongo.Client) ItemRepositoryService {
	return &itemRepository{db}
}

func (r *itemRepository) itemDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("item_db")
}
