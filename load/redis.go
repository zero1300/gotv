package load

import "github.com/go-redis/redis/v9"

func NewRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "172.17.0.1:6379",
	})
	return redisClient
}
