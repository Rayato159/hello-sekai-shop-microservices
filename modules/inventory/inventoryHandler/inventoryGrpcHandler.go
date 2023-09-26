package inventoryHandler

import "github.com/Rayato159/hello-sekai-shop-tutorial/modules/inventory/inventoryUsecase"

type (
	inventoryGrpcHandler struct {
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
	}
)

func NewInventoryGrpcHandler(inventoryUsecase inventoryUsecase.InventoryUsecaseService) *inventoryGrpcHandler {
	return &inventoryGrpcHandler{
		inventoryUsecase: inventoryUsecase,
	}
}
