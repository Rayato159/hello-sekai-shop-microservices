package inventoryHandler

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/inventory"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/inventory/inventoryUsecase"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/queue"
)

type (
	InventoryQueueHandlerService interface {
		AddPlayerItem()
		RemovePlayerItem()
		RollbackAddPlayerItem()
		RollbackRemovePlayerItem()
	}

	inventoryQueueHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
	}
)

func NewInventoryQueueHandler(cfg *config.Config, inventoryUsecase inventoryUsecase.InventoryUsecaseService) InventoryQueueHandlerService {
	return &inventoryQueueHandler{
		cfg:              cfg,
		inventoryUsecase: inventoryUsecase,
	}
}

func (h *inventoryQueueHandler) InventoryConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.inventoryUsecase.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("inventory", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("inventory", 0, 0)
		if err != nil {
			log.Println("Error: InventoryConsumer failed: ", err.Error())
			return nil, err
		}
	}

	return consumer, nil
}

func (h *inventoryQueueHandler) AddPlayerItem() {
	ctx := context.Background()

	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start AddPlayerItem ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: AddPlayerItem failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				h.inventoryUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(inventory.UpdateInventoryReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.AddPlayerItemRes(ctx, h.cfg, req)

				log.Printf("AddPlayerItem | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop AddPlayerItem...")
			return
		}
	}
}

func (h *inventoryQueueHandler) RollbackAddPlayerItem() {
	ctx := context.Background()

	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RollbackAddPlayerItem ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollbackAddPlayerItem failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "radd" {
				h.inventoryUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(inventory.RollbackPlayerInventoryReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.RollbackAddPlayerItem(ctx, h.cfg, req)

				log.Printf("RollbackAddPlayerItem | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop RollbackAddPlayerItem...")
			return
		}
	}
}

func (h *inventoryQueueHandler) RemovePlayerItem()         {}
func (h *inventoryQueueHandler) RollbackRemovePlayerItem() {}
