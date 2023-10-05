package authRepository

import (
	"context"

	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
	"github.com/Rayato159/hello-sekai-shop-tutorial/modules/auth"
	playerPb "github.com/Rayato159/hello-sekai-shop-tutorial/modules/player/playerPb"
	"github.com/Rayato159/hello-sekai-shop-tutorial/pkg/jwtauth"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func (m *AuthRepositoryMock) CredentialSearch(pctx context.Context, grpcUrl string, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error) {
	args := m.Called(pctx, grpcUrl, req)
	return args.Get(0).(*playerPb.PlayerProfile), args.Error(1)
}

func (m *AuthRepositoryMock) AccessToken(cfg *config.Config, claims *jwtauth.Claims) string {
	args := m.Called(cfg, claims)
	return args.String(0)
}

func (m *AuthRepositoryMock) RefreshToken(cfg *config.Config, claims *jwtauth.Claims) string {
	args := m.Called(cfg, claims)
	return args.String(0)
}

func (m *AuthRepositoryMock) InsertOnePlayerCredential(pctx context.Context, req *auth.Credential) (primitive.ObjectID, error) {
	args := m.Called(pctx, req)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *AuthRepositoryMock) FindOnePlayerCredential(pctx context.Context, credentialId string) (*auth.Credential, error) {
	args := m.Called(pctx, credentialId)
	return args.Get(0).(*auth.Credential), args.Error(1)
}

func (m *AuthRepositoryMock) FindOnePlayerProfileToRefresh(pctx context.Context, grpcUrl string, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error) {
	return nil, nil
}

func (m *AuthRepositoryMock) UpdateOnePlayerCredential(pctx context.Context, credentialId string, req *auth.UpdateRefreshTokenReq) error {
	return nil
}

func (m *AuthRepositoryMock) DeleteOnePlayerCredential(pctx context.Context, credentialId string) (int64, error) {
	return 0, nil
}

func (m *AuthRepositoryMock) FindOneAccessToken(pctx context.Context, accessToken string) (*auth.Credential, error) {
	return nil, nil
}

func (m *AuthRepositoryMock) RolesCount(pctx context.Context) (int64, error) {
	return 0, nil
}
