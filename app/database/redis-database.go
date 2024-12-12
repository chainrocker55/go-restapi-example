package database

import (
	"context"
	"creditlimit-connector/app/configs"
	"crypto/tls"
	"sync"

	"creditlimit-connector/app/log"

	"github.com/redis/go-redis/v9"
)

var (
	onceRedis   sync.Once
	redisClient *redis.Client
)

func InitRedisClient() *redis.Client {
	onceRedis.Do(func() {
		redisConf := configs.Conf.Database.RedisDatabase
		host := redisConf.Host + ":" + redisConf.Port
		pwd := redisConf.Password
		db := redisConf.DB
		var tlsConfig *tls.Config
		if redisConf.SSLMode {
			tlsConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
		}
		redisClient = redis.NewClient(
			&redis.Options{
				Addr:      host,
				Password:  pwd,
				DB:        db,
				TLSConfig: tlsConfig,
			},
		)

		if err := redisClient.Ping(context.TODO()).Err(); err != nil {
			log.Fatalf("Cannot connect redis database err: %v", err)
		}

		log.Info("Connected to redis database success")
	})
	return redisClient

}
