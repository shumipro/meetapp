package models

import (
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

type UserAuth struct {
	UserID        string `bson:"_id"      json:"userId"`
	TwitterToken  string `                json:"twitterToken"`
	FacebookToken string `                json:"facebookToken"`
}

type _UserAuthTable struct {
}

func (_ _UserAuthTable) Name() string {
	return "user_auth"
}

var _ modelsTable = (*_UserAuthTable)(nil)

var UserAuthTable = _UserAuthTable{}

func (t _UserAuthTable) withCollection(ctx context.Context, fn func(c *mgo.Collection)) {
	withDefaultCollection(ctx, t.Name(), fn)
}

// ----------------------------------------------

func (t _UserAuthTable) FindID(ctx context.Context, userID string) (result UserAuth, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.FindId(userID).One(&result)
	})
	return
}

func (t _UserAuthTable) Upsert(ctx context.Context, user UserAuth) error {
	var err error
	t.withCollection(ctx, func(c *mgo.Collection) {
		var result interface{} // bson.M
		_, err = c.FindId(user.UserID).Apply(mgo.Change{
			Update: user,
			Upsert: true,
		}, &result)
	})
	return err
}
