package authusecase

import (
  "awesomeProject/internal/entity"
  "context"
)

type userUseCase interface {
  RegisterUser(ctx context.Context, req entity.CreateUsrRequest) (token string, err error)
  LoginUser(ctx context.Context, req entity.LoginRequest) (token string, err error)
}

type UserUseCaseImpl struct {
  userRepository userUseCase
}

func NewUserUseCase(userRepository userUseCase) *UserUseCaseImpl {
  return &UserUseCaseImpl{userRepository: userRepository}
}

func (u *UserUseCaseImpl) RegisterUser(ctx context.Context, req entity.CreateUsrRequest) (token string, err error) {
  return u.userRepository.RegisterUser(ctx, req)
}

func (u *UserUseCaseImpl) LoginUser(ctx context.Context, req entity.LoginRequest) (token string, err error) {
  return u.userRepository.LoginUser(ctx, req)
}
