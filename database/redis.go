package database

import (
	"context"

	"github.com/2307vivek/song-lyrics/utils"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var ctx context.Context

func ConnectToRedis(url string) {
	opt, err := redis.ParseURL(url)
	utils.FailOnError(err, "Failed while connecting to redis.")

	Rdb = redis.NewClient(opt)

	createContext()
}

func CheckCache(key string, url string) bool {
	exists :=  Rdb.BFExists(ctx, key, url).Val()
	return exists
}

func AddToCache(key string, url string) bool {
	isadded := Rdb.BFAdd(ctx, key, url).Val()
	return isadded
}

func createContext() {
	ctx = context.Background()
}