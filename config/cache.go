package config

import (
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var (
	redisAddr     = os.Getenv("REDIS_ADDR")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisDB, _    = strconv.Atoi(os.Getenv("REDIS_DB"))
)

func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{Addr: redisAddr, Password: redisPassword, DB: redisDB})
	return client
}
