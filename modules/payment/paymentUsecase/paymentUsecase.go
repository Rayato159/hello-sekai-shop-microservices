package paymentUsecase

import (
	"context"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/inventory"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/item"
	itemPb "github.com/Rayato159/hello-sekai-shop-tutorial/modules/item/itemPb"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/payment"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/payment/paymentRepository"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/player"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/queue"
)

type (
	PaymentUsecaseService interface {
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
		FindItemsInIds(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error
		BuyItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error)
		SellItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error)
	}

	paymentUsecase struct {
		paymentRepository paymentRepository.PaymentRepositoryService
	}
)

func NewPaymentUsecase(paymentRepository paymentRepository.PaymentRepositoryService) PaymentUsecaseService {
	return &paymentUsecase{
		paymentRepository: paymentRepository,
	}
}

func (u *paymentUsecase) GetOffset(pctx context.Context) (int64, error) {
	return u.paymentRepository.GetOffset(pctx)
}
func (u *paymentUsecase) UpserOffset(pctx context.Context, offset int64) error {
	return u.paymentRepository.UpserOffset(pctx, offset)
}

func (u *paymentUsecase) PaymentConsumer(pctx context.Context, cfg *config.Config) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{cfg.Kafka.Url}, cfg.Kafka.ApiKey, cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := u.paymentRepository.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("payment", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("payment", 0, 0)
		if err != nil {
			log.Println("Error: PaymentConsumer failed: ", err.Error())
			return nil, err
		}
	}

	return consumer, nil
}

func (u *paymentUsecase) BuyOrSellConsumer(pctx context.Context, key string, cfg *config.Config, resCh chan<- *payment.PaymentTransferRes) {
	consumer, err := u.PaymentConsumer(pctx, cfg)
	if err != nil {
		resCh <- nil
		return
	}
	defer consumer.Close()

	log.Println("Start BuyOrSellConsumer ...")

	select {
	case err := <-consumer.Errors():
		log.Println("Error: BuyOrSellConsumer failed: ", err.Error())
		resCh <- nil
		return
	case msg := <-consumer.Messages():
		if string(msg.Key) == key {
			u.UpserOffset(pctx, msg.Offset+1)

			req := new(payment.PaymentTransferRes)

			if err := queue.DecodeMessage(req, msg.Value); err != nil {
				resCh <- nil
				return
			}

			resCh <- req
			log.Printf("BuyOrSellConsumer | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
		}
	}
}

func (u *paymentUsecase) BuyItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error) {
	if err := u.FindItemsInIds(pctx, cfg.Grpc.ItemUrl, req.Items); err != nil {
		return nil, err
	}

	stage1 := make([]*payment.PaymentTransferRes, 0)
	for _, item := range req.Items {
		u.paymentRepository.DockedPlayerMoney(pctx, cfg, &player.CreatePlayerTransactionReq{
			PlayerId: playerId,
			Amount:   -item.Price,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(pctx, "buy", cfg, resCh)

		res := <-resCh
		if res != nil {
			log.Println(res)
			stage1 = append(stage1, &payment.PaymentTransferRes{
				InventoryId:   "",
				TransactionId: res.TransactionId,
				PlayerId:      playerId,
				ItemId:        item.ItemId,
				Amount:        item.Price,
				Error:         res.Error,
			})
		}
	}

	for _, s1 := range stage1 {
		if s1.Error != "" {
			for _, ss1 := range stage1 {
				u.paymentRepository.RollbackTransaction(pctx, cfg, &player.RollbackPlayerTransactionReq{
					TransactionId: ss1.TransactionId,
				})
			}
			return nil, errors.New("error: buy item failed")
		}

	}

	stage2 := make([]*payment.PaymentTransferRes, 0)
	for _, s1 := range stage1 {
		u.paymentRepository.AddPlayerItem(pctx, cfg, &inventory.UpdateInventoryReq{
			PlayerId: playerId,
			ItemId:   s1.ItemId,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(pctx, "buy", cfg, resCh)

		res := <-resCh
		if res != nil {
			log.Println(res)
			stage2 = append(stage2, &payment.PaymentTransferRes{
				InventoryId:   res.InventoryId,
				TransactionId: s1.TransactionId,
				PlayerId:      playerId,
				ItemId:        s1.ItemId,
				Amount:        s1.Amount,
				Error:         s1.Error,
			})
		}
	}

	for _, s2 := range stage2 {
		if s2.Error != "" {
			for _, ss2 := range stage2 {
				u.paymentRepository.RollbackAddPlayerItem(pctx, cfg, &inventory.RollbackPlayerInventoryReq{
					InventoryId: ss2.InventoryId,
				})
			}

			for _, ss2 := range stage2 {
				u.paymentRepository.RollbackTransaction(pctx, cfg, &player.RollbackPlayerTransactionReq{
					TransactionId: ss2.TransactionId,
				})
			}

			return nil, errors.New("error: buy item failed")
		}
	}

	return stage2, nil
}

func (u *paymentUsecase) SellItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error) {
	if err := u.FindItemsInIds(pctx, cfg.Grpc.ItemUrl, req.Items); err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *paymentUsecase) FindItemsInIds(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error {
	setIds := make(map[string]bool)
	for _, v := range req {
		if !setIds[v.ItemId] {
			setIds[v.ItemId] = true
		}
	}

	itemData, err := u.paymentRepository.FindItemsInIds(pctx, grpcUrl, &itemPb.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0)
			for k := range setIds {
				itemIds = append(itemIds, k)
			}
			return itemIds
		}(),
	})
	if err != nil {
		log.Printf("Error: FindItemsInIds failed: %s", err.Error())
		return errors.New("error: items not found")
	}

	itemMaps := make(map[string]*item.ItemShowCase)
	for _, v := range itemData.Items {
		itemMaps[v.Id] = &item.ItemShowCase{
			ItemId:   v.Id,
			Title:    v.Title,
			Price:    v.Price,
			ImageUrl: v.ImageUrl,
			Damage:   int(v.Damage),
		}
	}

	for i := range req {
		if _, ok := itemMaps[req[i].ItemId]; !ok {
			log.Printf("Error: FindItemsInIds failed: %s", err.Error())
			return errors.New("error: items not found")
		}
		req[i].Price = itemMaps[req[i].ItemId].Price
	}

	return nil
}
