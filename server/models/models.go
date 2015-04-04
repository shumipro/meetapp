package models

import (
	"github.com/shumipro/meetapp/server/db"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

type modelsContext interface {
	withCollection(fn func(c *mgo.Collection))
}

func findAndModify(ctx modelsContext, findQuery bson.M, query bson.M) error {
	var result interface{}
	var err error
	ctx.withCollection(func(c *mgo.Collection) {
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
	m := db.MongoDB(ctx).Clone()
	defer m.Close()
	col := m.DB(db.Name()).C(name)
	fn(col)
}

func runTxCollection(ctx context.Context, ops []txn.Op) error {
	m := db.MongoDB(ctx).Clone()
	defer m.Close()
	col := m.DB(db.Name()).C("tx")

	runner := txn.NewRunner(col)
	return runner.Run(ops, "", nil)
}
