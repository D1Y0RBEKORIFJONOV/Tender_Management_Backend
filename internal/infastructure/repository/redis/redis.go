package redisrepository

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisUserRepository struct {
	redisClient *redis.Client
}

func NewRedis(cfg config.Config) *RedisUserRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return &RedisUserRepository{
		redisClient: client,
	}
}

func (r *RedisUserRepository) SaveUserReq(ctx context.Context, user entity.SaveRegisRequest, ttl time.Duration, key string) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key += fmt.Sprintf(":%s", user.Email)
	err = r.redisClient.Set(ctx, key, string(userJson), ttl).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r *RedisUserRepository) GetUserRegister(ctx context.Context, email, key string) (*entity.SaveRegisRequest, error) {
	key += fmt.Sprintf(":%s", email)

	val, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var user entity.SaveRegisRequest
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {

		return nil, err
	}

	return &user, nil
}

func (r *RedisUserRepository) SaveToCache(ctx context.Context, ttl time.Duration, key string, value []byte) error {
	err := r.redisClient.Set(ctx, key, string(value), ttl).Err()
	if err != nil {
		return err
	}
	return nil
}
func (r *RedisUserRepository) GetFromCache(ctx context.Context, key string) ([]byte, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}
func (r *RedisUserRepository) DeleteFromCache(ctx context.Context, key string) error {
	err := r.redisClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
