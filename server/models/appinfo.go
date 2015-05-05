package models

import (
	"time"

	"strings"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AppInfo アプリ
type AppInfo struct {
	ID            string               `bson:"_id" json:"id"`                 // アプリID
	Name          string               `           json:"name"`               // アプリ名
	Description   string               `           json:"description"`        // アプリ詳細
	Category      CategoryType         `           json:"category"`           // カテゴリ
	Platform      PlatformType         `           json:"platform"`           // プラットフォーム
	Language      LanguageType         `           json:"pLang"`              // プログラミング言語
	Keywords      string               `           json:"keywords"`           // フリーキーワード
	MainImage     string               `           json:"mainImageUrl"`       // メイン画像
	ImageURLs     []URLInfo            `           json:"images"`             // 紹介画像URLたち
	Area          AreaType             `           json:"meetingArea"`        // 場所
	StartDate     string               `           json:"projectStartDate"`   // 開始日
	ReleaseDate   string               `           json:"projectReleaseDate"` // リリース予定日
	GitHubURL     string               `           json:"githubUrl"`          // GitHubのURL
	DemoURL       string               `           json:"demoUrl"`            // デモURL
	Frequency     MeetingFrequencyType `           json:"meetingFrequency"`   // 頻度
	StarCount     int                  `           json:"starCount"`          // スター数
	Members       []Member             `           json:"currentMembers"`     // メンバー
	RecruitMember []RecruitInfo        `           json:"recruitMembers"`     // 募集メンバー
	Discussions   []DiscussionInfo     `           json:"discussions"`        // 「聞いてみる」の内容
	StarUsers     []string             `           json:"starUsers"`          // 「聞いてみる」の内容
	CreateAt      time.Time            `           json:"-"`
	UpdateAt      time.Time            `           json:"-"`
}

func (a AppInfo) FirstImageURL() string {
	if len(a.ImageURLs) > 0 {
		return a.ImageURLs[0].URL // TODO: とりあえず1件目をメインの画像にする
	} else {
		// set default image
		return "/img/no_img.png"
	}
}

func (a AppInfo) IsAdmin(userID string) bool {
	for _, m := range a.Members {
		if !m.IsAdmin {
			continue
		}

		if m.UserID == userID {
			return true
		}
	}
	return false
}

func (a AppInfo) Stared(userID string) bool {
	for _, a := range a.StarUsers {
		if strings.EqualFold(a, userID) {
			return true
		}
	}
	return false
}

func (d DiscussionInfo) FormatTime() string {
	return d.Timestamp.Format("2006-01-02 15:04")
}

// URLInfo 各種URL情報
type URLInfo struct {
	URL string `json:"url"`
}

type RecruitInfo struct {
	Occupation OccupationType `json:"occupation"` // 肩書とか役割
}

type Member struct {
	UserID     string         `json:"id"`
	Occupation OccupationType `json:"occupation"` // 肩書とか役割
	IsAdmin    bool           `json:"isAdmin"`    // 管理者フラグ
}

type DiscussionInfo struct {
	UserID    string    `json:"userId"`    // ユーザー
	Message   string    `json:"message"`   // コメント
	Timestamp time.Time `json:"timestamp"` // 投稿日時
}

// AppsContext appsのコレクション
type _AppsInfoTable struct {
}

func (_ _AppsInfoTable) Name() string {
	return "apps"
}

var _ modelsTable = (*_AppsInfoTable)(nil)

// AppsInfoTable appInfo
var AppsInfoTable = _AppsInfoTable{}

func (t _AppsInfoTable) withCollection(ctx context.Context, fn func(c *mgo.Collection)) {
	withDefaultCollection(ctx, t.Name(), fn)
}

func (t _AppsInfoTable) FindID(ctx context.Context, appID string) (result AppInfo, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.FindId(appID).One(&result)
	})
	return
}

func (t _AppsInfoTable) FindAll(ctx context.Context) (result []AppInfo, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.Find(bson.M{}).All(&result)
	})
	return
}

func (t _AppsInfoTable) FindFilter(ctx context.Context, filter AppInfoFilter, offset, num int) (maxLength int, result []AppInfo, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		query := c.Find(filter.Condition())

		maxLength, _ = query.Count()
		err = query.Skip(offset).Limit(num).All(&result)
	})
	return
}

func (t _AppsInfoTable) FindByAdminID(ctx context.Context, adminUserID string) (result []AppInfo, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.Find(bson.M{"members.userid": adminUserID, "members.isadmin": true}).All(&result)
	})
	return
}

func (t _AppsInfoTable) FindByJoinID(ctx context.Context, joinUserID string) (result []AppInfo, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.Find(bson.M{"members.userid": joinUserID}).Sort("createat").All(&result)
	})
	return
}

func (t _AppsInfoTable) FindLatest(ctx context.Context, offset int, num int) (result []AppInfo, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.Find(bson.M{}).Sort("createat").Skip(offset).Limit(num).All(&result)
	})
	return
}

func (t _AppsInfoTable) FindPopular(ctx context.Context, offset int, num int) (result []AppInfo, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.Find(bson.M{}).Sort("-starcount").Skip(offset).Limit(num).All(&result)
	})
	return
}

// Upsert 登録
func (t _AppsInfoTable) Upsert(ctx context.Context, app AppInfo) error {
	var err error
	t.withCollection(ctx, func(c *mgo.Collection) {
		var result interface{} // bson.M
		_, err = c.FindId(app.ID).Apply(mgo.Change{
			Update: app,
			Upsert: true,
		}, &result)
	})
	return err
}

func (t _AppsInfoTable) Delete(ctx context.Context, appID string) (err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.RemoveId(appID)
	})
	return
}

// document単位でatomicな更新
func (t _AppsInfoTable) findAndModify(ctx context.Context, findQuery bson.M, query bson.M) error {
	return findAndModify(t, ctx, findQuery, query)
}
