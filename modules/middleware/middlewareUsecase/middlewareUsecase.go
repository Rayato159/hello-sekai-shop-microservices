package middlewareUsecase

import (
	"errors"

	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/middleware/middlewareRepository"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/jwtauth"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/rbac"
	"github.com/labstack/echo/v4"
)

type (
	MiddlewareUsecaseService interface {
		JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error)
		RbacAuthorization(c echo.Context, cfg *config.Config, expected []int) (echo.Context, error)
	}

	middlewareUsecase struct {
		middlewareRepository middlewareRepository.MiddlewareRepositoryService
	}
)

func NewMiddlewareUsecase(middlewareRepository middlewareRepository.MiddlewareRepositoryService) MiddlewareUsecaseService {
	return &middlewareUsecase{middlewareRepository}
}

func (u *middlewareUsecase) JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error) {
	ctx := c.Request().Context()

	claims, err := jwtauth.ParseToken(cfg.Jwt.AccessSecretKey, accessToken)
	if err != nil {
		return nil, err
	}

	if err := u.middlewareRepository.AccessTokenSearch(ctx, cfg.Grpc.AuthUrl, accessToken); err != nil {
		return nil, err
	}

	c.Set("player_id", claims.PlayerId)
	c.Set("role_code", claims.RoleCode)

	return c, nil
}

func (u *middlewareUsecase) RbacAuthorization(c echo.Context, cfg *config.Config, expected []int) (echo.Context, error) {
	ctx := c.Request().Context()

	playerRoleCode := c.Get("role_code").(int)

	rolesCount, err := u.middlewareRepository.RolesCount(ctx, cfg.Grpc.AuthUrl)
	if err != nil {
		return nil, err
	}

	playerRoleBinary := rbac.IntToBinary(playerRoleCode, int(rolesCount))

	for i := 0; i < int(rolesCount); i++ {
		if playerRoleBinary[i]&expected[i] == 1 {
			return c, nil
		}
	}

	return nil, errors.New("error: permission denied")
}
