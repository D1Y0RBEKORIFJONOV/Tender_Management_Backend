package auth

import (
	"awesomeProject/internal/entity"
	redisrepository "awesomeProject/internal/infastructure/repository/redis"
	token2 "awesomeProject/internal/infastructure/token"
	authusecase "awesomeProject/internal/usecase/auth"
	notificationusecase "awesomeProject/internal/usecase/notification"
	"context"
	"errors"
	"fmt"
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

func (a *Auth) RegisterUser(ctx context.Context, req entity.CreateUsrRequest) (token string, err error) {
	const op = "auth.register"
	log := a.logger.With(
		slog.String("method", op))
	log.Info("start")
	defer log.Info("end")
	ok, err := a.auth.IsHaveUser(ctx, req.Email)
	if err != nil {
		log.Error("err", err.Error())
		return "", errors.New("Email already exists")
	}
	if ok {
		return "", errors.New("Email already exists")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("err", err.Error())
		return "", err
	}

	user, err := a.auth.SaveUser(ctx, &entity.CreateUsrRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: string(passwordHash),
		Role:     req.Role,
	})
	if err != nil {
		log.Error("err", err.Error())
		return "", err
	}

	err = a.notification.CreateNotification(ctx, &entity.CreateNotification{
		UserId: user.ID,
	})
	if err != nil {
		log.Error("err", err.Error())
		return "", err
	}

	token, _, err = token2.GenerateTokens(user)
	if err != nil {
		log.Error("err", err.Error())
		return "", err
	}

	return token, nil
}

func (a *Auth) LoginUser(ctx context.Context, req entity.LoginRequest) (token string, err error) {
	const op = "auth.login"
	log := a.logger.With(slog.String("method", op))
	log.Info("start")
	defer log.Info("end")

	user, err := a.auth.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error("err", req.Email)
		log.Error("err", err.Error())
		return "", errors.New("User not found")
	}
	log.Info(fmt.Sprintf("req: %v", req))
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Error("Failed to Login", err.Error())
		return "", errors.New("Invalid username or password")
	}
	log.Info("Successfully logged in")
	token, _, err = token2.GenerateTokens(user)
	if err != nil {
		log.Error("err", err.Error())
		return "", err
	}
	return token, nil
}
