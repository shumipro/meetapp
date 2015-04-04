package db

import (
	"fmt"

	"golang.org/x/net/context"
	mgo "gopkg.in/mgo.v2"
)

const mongoDBName = "default"

type mongodb string

func MongoDB(ctx context.Context) *mgo.Session {
	key := mongodb(mongoDBName)
	db, _ := ctx.Value(key).(*mgo.Session)
	return db
}

func Name() string {
	return mongoDBName
}

func OpenMongoDB(ctx context.Context, host string, port int) context.Context {
	url := fmt.Sprintf("%s:%d", host, port)
	sesh, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, mongodb(mongoDBName), sesh)
	return ctx
}

func CloseMongoDB(ctx context.Context) context.Context {
	sesh := MongoDB(ctx)
	if sesh == nil {
		fmt.Println("not found mongoDB")
	}
	sesh.Close()
	ctx = context.WithValue(ctx, mongodb(mongoDBName), nil)
	return ctx
}
