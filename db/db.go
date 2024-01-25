package db

import "github.com/redis/go-redis/v9"

func DatabaseConnect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis-13823.c322.us-east-1-2.ec2.cloud.redislabs.com:13823",
		Password: "IvzRlpj1Hz3pB0x4JkLEVWA6AZnZbtPU",
		DB:       0,
	})
}
