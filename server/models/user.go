package models

import (
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

type User struct {
	ID     string       `bson:"_id"` // UUID自動生成
	Name   string       // FacebookNameと同じ?別名？ TODO: いらないかも
	FBUser FacebookUser `bson:"facebook"` // Facebookのme情報
}

type FacebookUser struct {
	ID          string `json:"id"`
//	Email       string `json:"email"` // TODO: いろいろ怖いのでとりません
	Name        string `json:"name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	Locale      string `json:"locale"`
	Link        string `json:"link"`
	Verified    bool   `json:"verified"`
	Timezone    int    `json:"timezone"`
	UpdatedTime string `json:"updated_time"`
}

type _UsersTable struct {
}

func (_ _UsersTable) Name() string {
	return "users"
}

// TODO: あとで...
//var _ modelsContext = (*_UsersTable)(nil)

func UsersTable() _UsersTable {
	return _UsersTable{}
}

func (t _UsersTable) withCollection(ctx context.Context, fn func(c *mgo.Collection)) {
	withDefaultCollection(ctx, t.Name(), fn)
}

// ----------------------------------------------

func (t _UsersTable) FindID(ctx context.Context, userID string) (result User, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.FindId(userID).One(&result)
	})
	return
}

func (t _UsersTable) FindByFacebookID(ctx context.Context, facebookID string) (result User, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.Find(bson.M{"facebook.id": facebookID}).One(&result)
	})
	return
}

// db.users.find({"facebook.name": {$regex: '.*Yoko.*', $options: "i"}}, {});
func (t _UsersTable) FindByKeyword(ctx context.Context, keyword string) (results []User, err error) {
	regexWord :=  fmt.Sprintf(".*%s.*", keyword)
	fmt.Println("Keyword = ", regexWord)

	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.Find(bson.M{"facebook.name": bson.M{
			"$regex": regexWord,
			"$options": "i",
		}}).All(&results)
	})
	return
}

// Upsert 登録
func (t _UsersTable) Upsert(ctx context.Context, user User) error {
	var err error
	t.withCollection(ctx, func(c *mgo.Collection) {
		var result interface{} // bson.M
		_, err = c.FindId(user.ID).Apply(mgo.Change{
			Update: user,
			Upsert: true,
		}, &result)
	})
	return err
}
