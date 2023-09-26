package itemHandler

import "github.com/Rayato159/hello-sekai-shop-tutorial/modules/item/itemUsecase"

type (
	itemGrpcHandler struct {
		itemUsecase itemUsecase.ItemUsecaseService
	}
)

func NewItemGrpcHandler(itemUsecase itemUsecase.ItemUsecaseService) *itemGrpcHandler {
	return &itemGrpcHandler{
		itemUsecase: itemUsecase,
	}
}
