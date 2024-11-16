package auth

import (
	"awesomeProject/internal/entity"
	redisrepository "awesomeProject/internal/infastructure/repository/redis"
	token2 "awesomeProject/internal/infastructure/token"
	authusecase "awesomeProject/internal/usecase/auth"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"math/rand"
	"time"
)

type Auth struct {
	logger    *slog.Logger
	repoRedis *redisrepository.RedisUserRepository
	auth      *authusecase.AuthDbUseCaseImpl
}

func NewAuth(logger *slog.Logger,
	repoRedis *redisrepository.RedisUserRepository,
	auth *authusecase.AuthDbUseCaseImpl) *Auth {
	return &Auth{
		logger:    logger,
		repoRedis: repoRedis,
		auth:      auth,
	}
}

func (a *Auth) RegisterUser(ctx context.Context, req entity.CreateUsrRequest) (message string, err error) {
	const op = "auth.register"
	log := a.logger.With(
		slog.String("method", op))
	log.Info("start")
	defer log.Info("end")
	ok, err := a.auth.IsHaveUser(ctx, req.Email)
	if err != nil {
		log.Error("err", err.Error())
		return "", err
	}
	if ok {
		return "User already exists", nil
	}
	randInt := rand.Int() % 1000
	err = a.repoRedis.SaveUserReq(ctx, entity.SaveRegisRequest{
		Username:   req.Username,
		Email:      req.Email,
		Password:   req.Password,
		Role:       req.Role,
		SecretCode: fmt.Sprintf("%d", randInt),
	}, time.Minute*10, "user:registration")

	if err != nil {
		log.Error("err", err.Error())
		return "", err
	}
	log.Info("Successfully registered user")
	return "Check your email", nil
}

func (a *Auth) VerifyUser(ctx context.Context, req entity.VerifyUserRequest) (entity.User, error) {
	const op = "auth.verify"
	log := a.logger.With(
		slog.String("method", op))
	log.Info("start")
	defer log.Info("end")
	user, err := a.repoRedis.GetUserRegister(ctx, req.Email, "user:registration")
	if err != nil {
		log.Error("err", err.Error())
		return entity.User{}, err
	}
	if user == nil {
		return entity.User{}, errors.New("user not found")
	}
	if user.SecretCode != req.SecretCode {
		return entity.User{}, errors.New("invalid secret code")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("err", err.Error())
		return entity.User{}, err
	}

	savedUser, err := a.auth.SaveUser(ctx, &entity.CreateUsrRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: string(passwordHash),
		Role:     user.Role,
	})

	if err != nil {
		log.Error("err", err.Error())
		return entity.User{}, err
	}

	return *savedUser, nil
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
