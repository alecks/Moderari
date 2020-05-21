package db

import (
	"moderari/internal/config"

	"github.com/go-redis/redis"
)

// Client is the main Redis client.
var Client *redis.Client

// Nil lets us avoid importing the redis package.
const Nil = redis.Nil

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.C.DB.Address,
		Password: config.C.DB.Password,
		DB:       config.C.DB.Database,
	})

	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
