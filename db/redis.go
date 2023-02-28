package db

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	rdb     *redis.Client
	expires time.Duration
}

// Create a new redis client (redis on 6379)
func NewRedisClient(host string, db int, expires time.Duration) (rc *RedisClient, err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "", // no password set
		DB:       db,
	})

	pong, err := rdb.Ping().Result()
	// log message
	log.Printf("Redis ping: %s", pong)

	rc = &RedisClient{
		rdb:     rdb,
		expires: expires,
	}

	return
}

func (rc *RedisClient) Set(key string, value interface{}) (err error) {
	json, err := json.Marshal(value)

	if err != nil {
		// log error
		log.Printf("Error marshalling line: %s", err)
		return
	}

	err = rc.rdb.Set(key, json, rc.expires).Err()
	return
}

func (rc *RedisClient) Get(key string) (value string, err error) {
	value, err = rc.rdb.Get(key).Result()
	return
}

func (rc *RedisClient) Del(key string) (err error) {
	err = rc.rdb.Del(key).Err()
	return
}
