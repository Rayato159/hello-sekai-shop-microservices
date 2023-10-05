package whydoweneedtest

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/auth"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/auth/authRepository"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/auth/authUsecase"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/player"
	playerPb "github.com/Rayato159/hello-sekai-shop-tutorial/modules/player/playerPb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CredentialSearch
// Email or password is invalid
// Success

// InsertOnePlayerCredential
// Success

// FindOnePlayerCredential
// Credential not found
// Success

// Cases -> 3

type (
	testLogin struct {
		ctx      context.Context
		cfg      *config.Config
		req      *auth.PlayerLoginReq
		expected *auth.ProfileIntercepter
		isErr    bool
	}
)

func TestLogin(t *testing.T) {
	repoMock := new(authRepository.AuthRepositoryMock)
	usecase := authUsecase.NewAuthUsecase(repoMock)

	cfg := NewTestConfig()
	ctx := context.Background()

	credentialIdSuccess := primitive.NewObjectID()
	credentialIdFailed := primitive.NewObjectID()

	tests := []testLogin{
		{
			ctx: ctx,
			cfg: cfg,
			req: &auth.PlayerLoginReq{
				Email:    "success@sekai.com",
				Password: "123456",
			},
			expected: &auth.ProfileIntercepter{
				PlayerProfile: &player.PlayerProfile{
					Id:        "player:001",
					Email:     "success@sekai.com",
					Username:  "player001",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				Credential: &auth.CredentialRes{
					Id:           credentialIdSuccess.Hex(),
					PlayerId:     "player:001",
					RoleCode:     0,
					AccessToken:  "xxx",
					RefreshToken: "xxx",
					CreatedAt:    time.Time{},
					UpdatedAt:    time.Time{},
				},
			},
			isErr: false,
		},
		{
			ctx: ctx,
			cfg: cfg,
			req: &auth.PlayerLoginReq{
				Email:    "failed2@sekai.com",
				Password: "123456",
			},
			expected: nil,
			isErr:    true,
		},
		{
			ctx: ctx,
			cfg: cfg,
			req: &auth.PlayerLoginReq{
				Email:    "failed3@sekai.com",
				Password: "123456",
			},
			expected: nil,
			isErr:    true,
		},
	}

	// CredentialSearch
	repoMock.On("CredentialSearch", ctx, cfg.Grpc.PlayerUrl, &playerPb.CredentialSearchReq{
		Email:    "success@sekai.com",
		Password: "123456",
	}).Return(&playerPb.PlayerProfile{
		Id:        "001",
		Email:     "success@sekai.com",
		Username:  "player001",
		RoleCode:  0,
		CreatedAt: "0001-01-01 00:00:00 +0000 UTC",
		UpdatedAt: "0001-01-01 00:00:00 +0000 UTC",
	}, nil)

	repoMock.On("CredentialSearch", ctx, cfg.Grpc.PlayerUrl, &playerPb.CredentialSearchReq{
		Email:    "failed2@sekai.com",
		Password: "123456",
	}).Return(&playerPb.PlayerProfile{}, errors.New("error: email or password is invalid"))

	repoMock.On("CredentialSearch", ctx, cfg.Grpc.PlayerUrl, &playerPb.CredentialSearchReq{
		Email:    "failed3@sekai.com",
		Password: "123456",
	}).Return(&playerPb.PlayerProfile{
		Id:        "003",
		Email:     "failed3@sekai.com",
		Username:  "player003",
		RoleCode:  0,
		CreatedAt: "0001-01-01 00:00:00 +0000 UTC",
		UpdatedAt: "0001-01-01 00:00:00 +0000 UTC",
	}, nil)

	// Access Token
	repoMock.On("AccessToken", cfg, mock.AnythingOfType("*jwtauth.Claims")).Return("xxx")

	// Refresh Token
	repoMock.On("RefreshToken", cfg, mock.AnythingOfType("*jwtauth.Claims")).Return("xxx")

	// InsertOnePlayerCredential
	repoMock.On("InsertOnePlayerCredential", ctx, &auth.Credential{
		PlayerId:     "player:001",
		RoleCode:     0,
		AccessToken:  "xxx",
		RefreshToken: "xxx",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}).Return(credentialIdSuccess, nil)

	repoMock.On("InsertOnePlayerCredential", ctx, &auth.Credential{
		PlayerId:     "player:003",
		RoleCode:     0,
		AccessToken:  "xxx",
		RefreshToken: "xxx",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}).Return(credentialIdFailed, nil)

	// FindOnePlayerCredential
	repoMock.On("FindOnePlayerCredential", ctx, credentialIdSuccess.Hex()).Return(&auth.Credential{
		Id:           credentialIdSuccess,
		PlayerId:     "player:001",
		RoleCode:     0,
		AccessToken:  "xxx",
		RefreshToken: "xxx",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}, nil)

	repoMock.On("FindOnePlayerCredential", ctx, credentialIdFailed.Hex()).Return(&auth.Credential{}, errors.New("error: player credential not found"))

	for i, test := range tests {
		fmt.Printf("case -> %d\n", i+1)

		result, err := usecase.Login(test.ctx, test.cfg, test.req)

		if test.isErr {
			assert.NotEmpty(t, err)
		} else {
			result.CreatedAt = time.Time{}
			result.UpdatedAt = time.Time{}
			result.Credential.CreatedAt = time.Time{}
			result.Credential.UpdatedAt = time.Time{}

			assert.Equal(t, test.expected, result)
		}
	}
}
