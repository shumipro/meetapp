package models

import (
	"fmt"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID            string        `bson:"_id"      json:"ID"`            // UUID自動生成
	Name          string        `                json:"Name"`          // ユーザー名
	ImageName     string        `                json:"ImageName"`     // アップロードしたファイル名
	ImageURL      string        `                json:"ImageURL"`      // ユーザーアイコンのURL
	LargeImageURL string        `                json:"LargeImageURL"` // ユーザーアイコンの大きいURL
	Comment       string        `             	 json:"Comment"`        // ひとこと
	HomePageURL   string        `             	 json:"HomePageURL"`    // ウェブサイトURL
	GitHubURL     string        `             	 json:"GitHubURL"`      // Github URL
	FBUser        FacebookUser  `bson:"facebook" json:"FBUser"`        // Facebookのme情報
	TwitterUser   anaconda.User `bson:"twitter"  json:"TwitterUser"`   // Twitterのshows情報
	CreateAt      time.Time     `                json:"-"`
	UpdateAt      time.Time     `                json:"-"`
}

func (u User) IconImageURL() string {
	if u.ImageURL != "" {
		return u.ImageURL
	}

	if u.FBUser.ID != "" {
		return fmt.Sprintf("https://graph.facebook.com/%s/picture?type=square", u.FBUser.ID)
	}

	if u.TwitterUser.Id != 0 {
		return u.TwitterUser.ProfileImageUrlHttps
	}

	return "/img/no_img/no_img_1.png"
}

func (u User) IconLargeImageURL() string {
	if u.LargeImageURL != "" {
		return u.LargeImageURL
	}

	if u.FBUser.ID != "" {
		return fmt.Sprintf("https://graph.facebook.com/%s/picture?type=large", u.FBUser.ID)
	}

	if u.TwitterUser.Id != 0 {
		return u.TwitterUser.ProfileImageUrlHttps
	}

	return "/img/no_img/no_img_1.png"
}

func (u User) IsEmpty() bool {
	return u.ID == ""
}

type FacebookUser struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Gender      string `json:"gender"`
	Locale      string `json:"locale"`
	Link        string `json:"link"`
	Verified    bool   `json:"verified"`
	Timezone    int    `json:"timezone"`
	UpdatedTime string `json:"updatedTime"`
}

type _UsersTable struct {
}

func (_ _UsersTable) Name() string {
	return "users"
}

var _ modelsTable = (*_UsersTable)(nil)

var UsersTable = _UsersTable{}

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

func (t _UsersTable) FindByTwitterID(ctx context.Context, twitterID int64) (result User, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.Find(bson.M{"twitter.id": twitterID}).One(&result)
	})
	return
}

// db.users.find({"facebook.name": {$regex: '.*Yoko.*', $options: "i"}}, {});
func (t _UsersTable) FindByKeyword(ctx context.Context, keyword string) (results []User, err error) {
	regexWord := fmt.Sprintf(".*%s.*", keyword)
	fmt.Println("Keyword = ", regexWord)

	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.Find(bson.M{"name": bson.M{
			"$regex":   regexWord,
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
