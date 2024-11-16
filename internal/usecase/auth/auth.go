package authusecase

import (
	"awesomeProject/internal/entity"
	"context"
)

type userUseCase interface {
	RegisterUser(ctx context.Context, req entity.CreateUsrRequest) (message string, err error)
	VerifyUser(ctx context.Context, req entity.VerifyUserRequest) (entity.User, error)
	LoginUser(ctx context.Context, req entity.LoginRequest) (entity.LoginResponse, error)
}

type UserUseCaseImpl struct {
	userRepository userUseCase
}

func NewUserUseCase(userRepository userUseCase) *UserUseCaseImpl {
	return &UserUseCaseImpl{userRepository: userRepository}
}

func (u *UserUseCaseImpl) RegisterUser(ctx context.Context, req entity.CreateUsrRequest) (string, error) {
	return u.userRepository.RegisterUser(ctx, req)
}

func (u *UserUseCaseImpl) VerifyUser(ctx context.Context, req entity.VerifyUserRequest) (entity.User, error) {
	return u.userRepository.VerifyUser(ctx, req)
}

func (u *UserUseCaseImpl) LoginUser(ctx context.Context, req entity.LoginRequest) (entity.LoginResponse, error) {
	return u.userRepository.LoginUser(ctx, req)
}
