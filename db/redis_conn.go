package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisDbConn struct {
	RedisClient *redis.Client
}

func GetRedisConn() *redis.Client {
	fmt.Println("Go Redis Tutorial")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}

func (rc RedisDbConn) SetEntry(k string, v interface{}, d time.Duration) error {
	j, _ := json.Marshal(v)
	return rc.RedisClient.Set(k, j, 10*time.Second).Err()
}

func (rc RedisDbConn) GetEntry(k string) (string, error) {
	return rc.RedisClient.Get(k).Result()
}
