package paymentHandler

import (
	"context"
	"net/http"

	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/payment"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/payment/paymentUsecase"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/request"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	PaymentHttpHandlerService interface {
		BuyItem(c echo.Context) error
		SellItem(c echo.Context) error
	}

	paymentHttpHandler struct {
		cfg            *config.Config
		paymentUsecase paymentUsecase.PaymentUsecaseService
	}
)

func NewPaymentHttpHandler(cfg *config.Config, paymentUsecase paymentUsecase.PaymentUsecaseService) PaymentHttpHandlerService {
	return &paymentHttpHandler{
		cfg:            cfg,
		paymentUsecase: paymentUsecase,
	}
}

func (h *paymentHttpHandler) BuyItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	playerId := c.Get("player_id").(string)

	req := &payment.ItemServiceReq{
		Items: make([]*payment.ItemServiceReqDatum, 0),
	}

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.paymentUsecase.BuyItem(ctx, h.cfg, playerId, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *paymentHttpHandler) SellItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	playerId := c.Get("player_id").(string)

	req := &payment.ItemServiceReq{
		Items: make([]*payment.ItemServiceReqDatum, 0),
	}

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.paymentUsecase.SellItem(ctx, h.cfg, playerId, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
