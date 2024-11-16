package authusecase

import (
	"awesomeProject/internal/entity"
	"context"
)

type authDbUseCase interface {
	SaveUser(ctx context.Context, req *entity.User) (*entity.User, error)
}

type AuthDbUseCaseImpl struct {
	authDbRepository authDbUseCase
}

func NewAuthDbUseCase(authDbRepository authDbUseCase) *AuthDbUseCaseImpl {
	return &AuthDbUseCaseImpl{authDbRepository: authDbRepository}
}

func (a *AuthDbUseCaseImpl) SaveUser(ctx context.Context, req *entity.CreateUsrRequest) (*entity.User, error) {
	return a.authDbRepository.SaveUser(ctx, req)
}

func (a *AuthDbUseCaseImpl) IsHaveUser(ctx context.Context, byEmail string) (bool, error) {
	return a.authDbRepository.IsHaveUser(ctx, byEmail)
}

func (a *AuthDbUseCaseImpl) GetUserByEmail(ctx context.Context, byEmail string) (*entity.User, error) {
	return a.authDbRepository.GetUserByEmail(ctx, byEmail)
}
