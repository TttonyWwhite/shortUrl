package models

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := Client.Ping().Result()
	fmt.Println(pong, err)
	if err == nil {
		log.Println("redis connection initialized")
	}
}

func AddLongToShort(longUrl string, shortUrl string) {
	// 缓存中数据过期时间暂且定为5分钟
	Client.Set(longUrl, shortUrl, time.Minute*5)
}

func AddShortToLong(shortUrl string, longUrl string) {
	Client.Set(shortUrl, longUrl, time.Minute*5)
}

func GetShortUrlFromRedis(longUrl string) (shortUrl string, err error) {
	shortUrl, err = Client.Get(longUrl).Result()
	if err != nil {
		return "", errors.New("url not in redis")
	}

	return shortUrl, nil
}

func GetLongUrlFromRedis(shortUrl string) (longUrl string, err error) {
	longUrl, err = Client.Get(shortUrl).Result()
	if err != nil {
		return "", errors.New("url not in redis")
	}

	return longUrl, nil
}
