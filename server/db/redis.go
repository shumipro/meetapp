package db

import (
	"fmt"
	"log"

	"os"

	"golang.org/x/net/context"
	"gopkg.in/redis.v2"
)

type redisDB string

func Redis(ctx context.Context) *redis.Client {
	key := redisDB("default")
	db, _ := ctx.Value(key).(*redis.Client)
	return db
}

func OpenRedis(ctx context.Context) context.Context {
	url := os.Getenv("REDISTOGO_URL")
	if url == "" {
		url = fmt.Sprintf("%s:%d", "localhost", 6379)
	}
	fmt.Println("redis", url)

	client := redis.NewTCPClient(&redis.Options{
		Addr: url,
	})
	ctx = context.WithValue(ctx, redisDB("default"), client)
	return ctx
}

func CloseRedis(ctx context.Context) context.Context {
	client := Redis(ctx)
	if client == nil {
		return ctx
	}

	if err := client.Close(); err != nil {
		log.Println("redis close error:", err)
	}

	return context.WithValue(ctx, redisDB("default"), nil)
}
