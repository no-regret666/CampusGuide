package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func InitRedis() (*redis.Client, error) {
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// 检查与Redis的连接是否成功
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("无法连接到Redis: %v", err)
	}
	fmt.Println("已连接到Redis", pong)

	return client, nil
}
