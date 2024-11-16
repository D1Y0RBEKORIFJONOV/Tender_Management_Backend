package authusecase

import (
	"awesomeProject/internal/entity"
	"context"
)

type authDbUseCase interface {
	SaveUser(ctx context.Context, req *entity.CreateUsrRequest) (*entity.User, error)
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
