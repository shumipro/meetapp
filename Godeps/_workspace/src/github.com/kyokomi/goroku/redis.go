package goroku

import (
	"fmt"
	"log"
	"net/url"
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
	addr, password := GetHerokuRedisAddr()
	client := redis.NewTCPClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	ctx = context.WithValue(ctx, redisDB("default"), client)
	return ctx
}

func GetHerokuRedisAddr() (addr string, password string) {
	addr = fmt.Sprintf("%s:%d", "localhost", 6379)
	password = ""

	redisURL := os.Getenv("REDISTOGO_URL")
	if redisURL == "" {
		fmt.Println("local: redis", addr, password)
		return
	}

	redisInfo, err := url.Parse(redisURL)
	if err != nil {
		return
	}

	addr = redisInfo.Host
	if redisInfo.User != nil {
		password, _ = redisInfo.User.Password()
	}
	return
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
