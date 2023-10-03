package paymentUsecase

import (
	"context"
	"errors"
	"log"

	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/item"
	itemPb "github.com/Rayato159/hello-sekai-shop-tutorial/modules/item/itemPb"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/payment"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/payment/paymentRepository"
)

type (
	PaymentUsecaseService interface {
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
		FindItemsInIds(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error
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
