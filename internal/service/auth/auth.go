package auth

import (
	"awesomeProject/internal/entity"
	redisrepository "awesomeProject/internal/infastructure/repository/redis"
	"context"
	"log/slog"
	"time"
)

type Auth struct {
	logger    *slog.Logger
	repoRedis redisrepository.RedisUserRepository
}

func (a *Auth) RegisterUser(ctx context.Context, req entity.CreateUsrRequest) (message string, err error) {
	a.repoRedis.SaveUserReq(ctx, req, time.Hour*2, "12344")
	panic("not implemented")
}
