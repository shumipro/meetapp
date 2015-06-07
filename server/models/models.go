package models

import (
	"github.com/kyokomi/goroku"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

type modelsTable interface {
	withCollection(ctx context.Context, fn func(c *mgo.Collection))
}

func findAndModify(t modelsTable, ctx context.Context, findQuery bson.M, query bson.M) error {
	var result interface{}
	var err error
	t.withCollection(ctx, func(c *mgo.Collection) {
		_, err = c.Find(findQuery).Apply(mgo.Change{
			Update: query,
		}, &result)
	})
	if err != nil {
		return err
	}
	return nil
}

func withDefaultCollection(ctx context.Context, name string, fn func(c *mgo.Collection)) {
	m := goroku.MustMongoDB(ctx).Clone()
	defer m.Close()
	col := m.DB(goroku.MongoDBName()).C(name)
	fn(col)
}

func runTxCollection(ctx context.Context, ops []txn.Op) error {
	m := goroku.MustMongoDB(ctx).Clone()
	defer m.Close()
	col := m.DB(goroku.MongoDBName()).C("tx")

	runner := txn.NewRunner(col)
	return runner.Run(ops, "", nil)
}
