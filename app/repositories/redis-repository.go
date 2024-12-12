package repositories

import (
	"context"
	"creditlimit-connector/app/database"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	onceRedisRepo sync.Once
	redisRepo     RedisRepository
)

type RedisRepository interface {
	Save(key string, value interface{}, ttlInSecond int) error
	Find(key string, model interface{}) error
}

type RedisRepositoryImp struct {
	redisClient *redis.Client
}

func InitRedisRepository() RedisRepository {
	onceRedisRepo.Do(func() {
		redisRepo = &RedisRepositoryImp{
			redisClient: database.InitRedisClient(),
		}
	})
	return redisRepo
}

func (r *RedisRepositoryImp) Save(key string, value interface{}, ttlInSecond int) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ctx := context.Background()
	ttl := time.Duration(ttlInSecond) * time.Second
	err = r.redisClient.Set(ctx, key, v, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepositoryImp) Find(key string, model interface{}) error {
	ctx := context.Background()
	val, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), model)
	return err
}
