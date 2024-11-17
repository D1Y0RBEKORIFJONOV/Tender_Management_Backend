package auth

import (
	"awesomeProject/internal/entity"
	redisrepository "awesomeProject/internal/infastructure/repository/redis"
	token2 "awesomeProject/internal/infastructure/token"
	authusecase "awesomeProject/internal/usecase/auth"
	notificationusecase "awesomeProject/internal/usecase/notification"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type Auth struct {
	logger       *slog.Logger
	repoRedis    *redisrepository.RedisUserRepository
	auth         *authusecase.AuthDbUseCaseImpl
	notification *notificationusecase.NotificationUseCase
}

func NewAuth(logger *slog.Logger,
	repoRedis *redisrepository.RedisUserRepository,
	auth *authusecase.AuthDbUseCaseImpl,
	notification *notificationusecase.NotificationUseCase) *Auth {
	return &Auth{
		logger:       logger,
		repoRedis:    repoRedis,
		auth:         auth,
		notification: notification,
	}
}

func (a *Auth) RegisterUser(ctx context.Context, req entity.CreateUsrRequest) (entity.User, error) {
	const op = "auth.register"
	log := a.logger.With(
		slog.String("method", op))
	log.Info("start")
	defer log.Info("end")
	ok, err := a.auth.IsHaveUser(ctx, req.Email)
	if err != nil {
		log.Error("err", err.Error())
		return entity.User{}, err
	}
	if ok {
		return entity.User{}, nil
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("err", err.Error())
		return entity.User{}, err
	}

	user, err := a.auth.SaveUser(ctx, &entity.CreateUsrRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: string(passwordHash),
		Role:     req.Role,
	})
	if err != nil {
		log.Error("err", err.Error())
		return entity.User{}, err
	}

	err = a.notification.CreateNotification(ctx, &entity.CreateNotification{
		UserId: user.ID,
	})
	if err != nil {
		log.Error("err", err.Error())
		return entity.User{}, err
	}

	return *user, nil
}

func (a *Auth) LoginUser(ctx context.Context, req entity.LoginRequest) (entity.LoginResponse, error) {
	const op = "auth.login"
	log := a.logger.With(slog.String("method", op))
	log.Info("start")
	defer log.Info("end")

	user, err := a.auth.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error("err", err.Error())
		return entity.LoginResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Error("Failed to login", err.Error())
		return entity.LoginResponse{}, errors.New("invalid password or username")
	}
	log.Info("Successfully logged in")
	token, _, err := token2.GenerateTokens(user)
	if err != nil {
		log.Error("err", err.Error())
		return entity.LoginResponse{}, err
	}
	return entity.LoginResponse{Token: token}, nil
}
