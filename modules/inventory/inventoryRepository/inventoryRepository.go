package inventoryRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	InventoryRepositoryService interface{}

	inventoryRepository struct {
		db *mongo.Client
	}
)

func NewInventoryRepository(db *mongo.Client) InventoryRepositoryService {
	return &inventoryRepository{db}
}

func (r *inventoryRepository) inventoryDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("inventory_db")
}
