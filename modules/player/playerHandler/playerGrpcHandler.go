package playerHandler

import "github.com/Rayato159/hello-sekai-shop-tutorial/modules/player/playerUsecase"

type (
	playerGrpcHandlerService struct {
		playerUsecase playerUsecase.PlayerUsecaseService
	}
)

func NewPlayerGrpcHandler(playerUsecase playerUsecase.PlayerUsecaseService) *playerGrpcHandlerService {
	return &playerGrpcHandlerService{playerUsecase}
}
