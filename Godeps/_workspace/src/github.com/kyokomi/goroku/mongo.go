package goroku

import (
	"fmt"

	"os"

	"net/url"

	"strings"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

const mongoDBName = "meetapp"

type mongodb string

var databaseName string

func MustMongoDB(ctx context.Context) (*mgo.Session) {
	db, ok := MongoDB(ctx)
	if !ok {
		panic("not found mongoDB")
	}
	return db
}

func MongoDB(ctx context.Context) (*mgo.Session, bool) {
	key := mongodb(mongoDBName)
	session, ok := ctx.Value(key).(*mgo.Session)
	return session, ok
}

func MongoDBName() string {
	return databaseName
}

func WithMockMongoDB() {
	databaseName = "test_" + mongoDBName
}

func OpenMongoDB(ctx context.Context) context.Context {
	uri, dbName := getHerokuMongoURI()
	databaseName = dbName

	sesh, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, mongodb(mongoDBName), sesh)
	return ctx
}

func getHerokuMongoURI() (uri string, dbName string) {
	// default
	uri = fmt.Sprintf("%s:%d", "localhost", 27017)
	dbName = mongoDBName

	mongoURI := os.Getenv("MONGOLAB_URI")
	if mongoURI == "" {
		fmt.Println("local: mongoDB", uri, dbName)
		return
	}
	mongoInfo, err := url.Parse(mongoURI)
	if err != nil {
		return
	}

	uri = mongoURI
	dbName = strings.Replace(mongoInfo.Path, "/", "", 1)
	return
}

func CloseMongoDB(ctx context.Context) context.Context {
	sesh, _ := MongoDB(ctx)
	if sesh == nil {
		fmt.Println("not found mongoDB")
	}
	sesh.Close()
	ctx = context.WithValue(ctx, mongodb(mongoDBName), nil)
	return ctx
}
