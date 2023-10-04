package paymentRepository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/inventory"
	itemPb "github.com/Rayato159/hello-sekai-shop-tutorial/modules/item/itemPb"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/models"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/player"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/grpccon"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/jwtauth"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/queue"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	PaymentRepositoryService interface {
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
		FindItemsInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
		DockedPlayerMoney(pctx context.Context, cfg *config.Config, req *player.CreatePlayerTransactionReq) error
		RollbackTransaction(pctx context.Context, cfg *config.Config, req *player.RollbackPlayerTransactionReq) error
		AddPlayerItem(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) error
		RollbackAddPlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq) error
		RemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) error
		RollbackRemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq) error
	}

	paymentRepository struct {
		db *mongo.Client
	}
)

func NewPaymentRepository(db *mongo.Client) PaymentRepositoryService {
	return &paymentRepository{db}
}

func (r *paymentRepository) paymentDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("payment_db")
}

func (r *paymentRepository) GetOffset(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payment_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: GetOffset failed: %s", err.Error())
		return -1, errors.New("error: GetOffset failed")
	}

	return result.Offset, nil
}

func (r *paymentRepository) UpserOffset(pctx context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.paymentDbConn(ctx)
	col := db.Collection("payment_queue")

	result, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error: UpserOffset failed: %s", err.Error())
		return errors.New("error: UpserOffset failed")
	}
	log.Printf("Info: UpserOffset result: %v", result)

	return nil
}

func (r *paymentRepository) FindItemsInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
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

func (r *paymentRepository) DockedPlayerMoney(pctx context.Context, cfg *config.Config, req *player.CreatePlayerTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: DockedPlayerMoney failed: %s", err.Error())
		return errors.New("error: docked player money failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"player",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: DockedPlayerMoney failed: %s", err.Error())
		return errors.New("error: docked player money failed")
	}

	return nil
}

func (r *paymentRepository) RollbackTransaction(pctx context.Context, cfg *config.Config, req *player.RollbackPlayerTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: DockedPlayerMoney failed: %s", err.Error())
		return errors.New("error: rollback player transaction failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"player",
		"rtransaction",
		reqInBytes,
	); err != nil {
		log.Printf("Error: DockedPlayerMoney failed: %s", err.Error())
		return errors.New("error: rollback player transaction failed")
	}

	return nil
}

func (r *paymentRepository) AddPlayerItem(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: AddPlayerItem failed: %s", err.Error())
		return errors.New("error: add player item failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: AddPlayerItem failed: %s", err.Error())
		return errors.New("error: add player item failed")
	}

	return nil
}

func (r *paymentRepository) RollbackAddPlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: RollbackAddPlayerItem failed: %s", err.Error())
		return errors.New("error: rollback add player item failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"radd",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RollbackAddPlayerItem failed: %s", err.Error())
		return errors.New("error: rollback add player item failed")
	}

	return nil
}

func (r *paymentRepository) RemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: RemovePlayerItem failed: %s", err.Error())
		return errors.New("error: remove player item failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"sell",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RemovePlayerItem failed: %s", err.Error())
		return errors.New("error: remove player item failed")
	}

	return nil
}

func (r *paymentRepository) RollbackRemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: RollbackRemovePlayerItem failed: %s", err.Error())
		return errors.New("error: rollback remove player item failed")
	}

	if err := queue.PushMessageWithKeyToQueue(
		[]string{cfg.Kafka.Url},
		cfg.Kafka.ApiKey,
		cfg.Kafka.Secret,
		"inventory",
		"rremove",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RollbackRemovePlayerItem failed: %s", err.Error())
		return errors.New("error: rollback remove player item failed")
	}

	return nil
}
