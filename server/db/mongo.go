package db

import (
	"fmt"

	"os"

	"golang.org/x/net/context"
	mgo "gopkg.in/mgo.v2"
	"strings"
)

const mongoDBName = "meetapp"

type mongodb string

var databaseName string

func MongoDB(ctx context.Context) *mgo.Session {
	key := mongodb(mongoDBName)
	db, _ := ctx.Value(key).(*mgo.Session)
	return db
}

func DBName() string {
	return databaseName
}

func OpenMongoDB(ctx context.Context) context.Context {
	url := os.Getenv("MONGOLAB_URI")
	if url == "" {
		url = fmt.Sprintf("%s:%d", "localhost", 27017)
	} else {
		// mongodb://<dbuser>:<dbpassword>@ds061371.mongolab.com:61371/heroku_app35413694st
		databaseName = strings.TrimLeft(url, "/")
		fmt.Println(databaseName)
	}
	fmt.Println("mongoDB", url)
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
