package playerHandler

import (
	"context"
	"net/http"
	"strings"

	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/player"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/player/playerUsecase"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/request"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	PlayerHttpHandlerService interface {
		CreatePlayer(c echo.Context) error
		FindOnePlayerProfile(c echo.Context) error
	}

	playerHttpHandler struct {
		cfg           *config.Config
		playerUsecase playerUsecase.PlayerUsecaseService
	}
)

func NewPlayerHttpHandler(cfg *config.Config, playerUsecase playerUsecase.PlayerUsecaseService) PlayerHttpHandlerService {
	return &playerHttpHandler{playerUsecase: playerUsecase}
}

func (h *playerHttpHandler) CreatePlayer(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(player.CreatePlayerReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.playerUsecase.CreatePlayer(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *playerHttpHandler) FindOnePlayerProfile(c echo.Context) error {
	ctx := context.Background()

	playerId := strings.TrimPrefix(c.Param("player_id"), "player:")

	res, err := h.playerUsecase.FindOnePlayerProfile(ctx, playerId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusBadRequest, res)
}
