package main

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.0.6:6379",
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})
	ctx := context.Background()
	s := rdb.Ping(ctx).String()
	val, _ := rdb.Get(ctx, "11").Result()
	print(val, s)
}
