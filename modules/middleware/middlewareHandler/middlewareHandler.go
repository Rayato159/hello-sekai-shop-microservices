package middlewareHandler

import (
	"net/http"
	"strings"

	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/middleware/middlewareUsecase"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	MiddlewareHandlerService interface {
		JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc
		RbacAuthorization(next echo.HandlerFunc, expected []int) echo.HandlerFunc
		PlayerIdParamValidation(next echo.HandlerFunc) echo.HandlerFunc
	}

	middlewareHandler struct {
		cfg               *config.Config
		middlewareUsecase middlewareUsecase.MiddlewareUsecaseService
	}
)

func NewMiddlewareHandler(cfg *config.Config, middlewareUsecase middlewareUsecase.MiddlewareUsecaseService) MiddlewareHandlerService {
	return &middlewareHandler{
		cfg:               cfg,
		middlewareUsecase: middlewareUsecase,
	}
}

func (h *middlewareHandler) JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")

		newCtx, err := h.middlewareUsecase.JwtAuthorization(c, h.cfg, accessToken)
		if err != nil {
			return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
		}

		return next(newCtx)
	}
}

func (h *middlewareHandler) RbacAuthorization(next echo.HandlerFunc, expected []int) echo.HandlerFunc {
	return func(c echo.Context) error {
		newCtx, err := h.middlewareUsecase.RbacAuthorization(c, h.cfg, expected)
		if err != nil {
			return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
		}

		return next(newCtx)
	}
}

func (h *middlewareHandler) PlayerIdParamValidation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		newCtx, err := h.middlewareUsecase.PlayerIdParamValidation(c)
		if err != nil {
			return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
		}

		return next(newCtx)
	}
}
