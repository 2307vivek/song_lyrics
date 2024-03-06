package database

import (
	"github.com/2307vivek/song-lyrics/utils"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func ConnectToRedis(url string) {
	opt, err := redis.ParseURL(url)
	utils.FailOnError(err, "Failed while connecting to redis.")

	Rdb = redis.NewClient(opt)
}