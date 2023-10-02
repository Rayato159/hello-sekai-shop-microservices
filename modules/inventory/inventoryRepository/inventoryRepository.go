package inventoryRepository

import (
	"context"
	"errors"
	"log"
	"time"

	itemPb "github.com/Rayato159/hello-sekai-shop-tutorial/modules/item/itemPb"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/grpccon"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/jwtauth"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	InventoryRepositoryService interface {
		FindItemsInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
	}

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

func (r *inventoryRepository) FindItemsInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	jwtauth.SetApiKeyInContext(&ctx)
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gRPC connection failed: %s", err.Error())
		return nil, errors.New("error: gRPC connection failed")
	}

	result, err := conn.Item().FindItemsInIds(ctx, req)
	if err != nil {
		log.Printf("Error: FindItemsInIds failed: %s", err.Error())
		return nil, errors.New("error: items not found")
	}

	if result == nil {
		log.Printf("Error: FindItemsInIds failed: %s", err.Error())
		return nil, errors.New("error: items not found")
	}

	if len(result.Items) == 0 {
		log.Printf("Error: FindItemsInIds failed: %s", err.Error())
		return nil, errors.New("error: items not found")
	}

	return result, nil
}
