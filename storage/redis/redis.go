package redis

import (
	"context"
	"log"

	"blog/config"

	"github.com/redis/go-redis/v9"
)

var Client *redis.ClusterClient

func Init() error {
	Client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    config.GlobalConfig.Redis.Cluster,
		Password: config.GlobalConfig.Redis.Password,
	})

	// 测试连接
	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("Redis cluster connected successfully")
	return nil
}
