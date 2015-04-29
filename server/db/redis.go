package db

import (
	"fmt"
	"log"

	"os"

	"net/url"

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
	redisURL := os.Getenv("REDISTOGO_URL")
	server := ""
	password := ""
	if redisURL == "" {
		server = fmt.Sprintf("%s:%d", "localhost", 6379)
	} else {
		redisInfo, _ := url.Parse(redisURL)
		server = redisInfo.Host
		if redisInfo.User != nil {
			password, _ = redisInfo.User.Password()
		}
	}
	fmt.Println("redis", server, password)

	client := redis.NewTCPClient(&redis.Options{
		Addr: server,
		Password: password,
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
