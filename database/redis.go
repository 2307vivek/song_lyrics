package database

import (
	"context"

	"github.com/2307vivek/song-lyrics/utils"
	"github.com/2307vivek/song-lyrics/utils/api"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var ctx context.Context

func ConnectToRedis(url string) {
	opt, err := redis.ParseURL(url)
	utils.FailOnError(err, "Failed while connecting to redis.")

	Rdb = redis.NewClient(opt)

	createContext()

	api.AppStatus.Connections.Redis = true
}

func ExistsInCache(cache string, item string) bool {
	exists :=  Rdb.BFExists(ctx, cache, item).Val()
	return exists
}

func AddToCache(cache string, item string) bool {
	isadded := Rdb.BFAdd(ctx, cache, item).Val()
	return isadded
}

func createContext() {
	ctx = context.Background()
}