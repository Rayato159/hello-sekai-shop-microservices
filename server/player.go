package server

import (
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/player/playerHandler"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/player/playerRepository"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/player/playerUsecase"
)

func (s *server) playerService() {
	repo := playerRepository.NewPlayerRepository(s.db)
	usecase := playerUsecase.NewPlayerUsecase(repo)
	httpHandler := playerHandler.NewPlayerHttpHandler(s.cfg, usecase)
	grpcHandler := playerHandler.NewPlayerGrpcHandler(usecase)
	queueHandler := playerHandler.NewPlayerQueueHandler(s.cfg, usecase)

	_ = httpHandler
	_ = grpcHandler
	_ = queueHandler

	player := s.app.Group("/player_v1")

	// Health Check
	_ = player
}
