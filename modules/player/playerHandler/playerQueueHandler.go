package playerHandler

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/player"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/player/playerUsecase"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/queue"
)

type (
	PlayerQueueHandlerService interface {
		DockedPlayerMoney()
		AddPlayerMoney()
		RollbackPlayerTransaction()
	}

	playerQueueHandler struct {
		cfg           *config.Config
		playerUsecase playerUsecase.PlayerUsecaseService
	}
)

func NewPlayerQueueHandler(cfg *config.Config, playerUsecase playerUsecase.PlayerUsecaseService) PlayerQueueHandlerService {
	return &playerQueueHandler{
		cfg:           cfg,
		playerUsecase: playerUsecase,
	}
}

func (h *playerQueueHandler) PlayerConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.playerUsecase.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("player", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("player", 0, 0)
		if err != nil {
			log.Println("Error: PaymentConsumer failed: ", err.Error())
			return nil, err
		}
	}

	return consumer, nil
}

func (h *playerQueueHandler) DockedPlayerMoney() {
	ctx := context.Background()

	consumer, err := h.PlayerConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start DockedPlayerMoney ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: DockedPlayerMoney failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				h.playerUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(player.CreatePlayerTransactionReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.playerUsecase.DockedPlayerMoneyRes(ctx, h.cfg, req)

				log.Printf("DockedPlayerMoney | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop DockedPlayerMoney...")
			return
		}
	}
}

func (h *playerQueueHandler) AddPlayerMoney() {
	ctx := context.Background()

	consumer, err := h.PlayerConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start AddPlayerMoney ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: AddPlayerMoney failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "sell" {
				h.playerUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(player.CreatePlayerTransactionReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.playerUsecase.AddPlayerMoneyRes(ctx, h.cfg, req)

				log.Printf("AddPlayerMoney | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop AddPlayerMoney...")
			return
		}
	}
}

func (h *playerQueueHandler) RollbackPlayerTransaction() {
	ctx := context.Background()

	consumer, err := h.PlayerConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RollbackPlayerTransaction ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollbackPlayerTransaction failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "rtransaction" {
				h.playerUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(player.RollbackPlayerTransactionReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.playerUsecase.RollbackPlayerTransaction(ctx, req)

				log.Printf("RollbackPlayerTransaction | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop RollbackPlayerTransaction...")
			return
		}
	}
}
