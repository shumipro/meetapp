package models

import (
	"time"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AppInfo アプリ
type AppInfo struct {
	ID            string           `bson:"_id" json:"id"`                 // アプリID
	Name          string           `           json:"name"`               // アプリ名
	Description   string           `           json:"description"`        // アプリ詳細
	Category      string           `           json:"category"`           // カテゴリ
	Platform      string           `           json:"platform"`           // プラットフォーム
	Language      string           `           json:"pLang"`              // プログラミング言語
	Keywords      string           `           json:"keywords"`           // フリーキーワード
	MainImage     string           `           json:"mainImageUrl"`       // メイン画像
	ImageURLs     []URLInfo        `           json:"images"`             // 紹介画像URLたち
	Area          string           `           json:"meetingArea"`        // 場所
	StartDate     string           `           json:"projectStartDate"`   // 開始日
	ReleaseDate   string           `           json:"projectReleaseDate"` // リリース予定日
	GitHubURL     string           `           json:"githubUrl"`          // GitHubのURL
	DemoURL       string           `           json:"demoUrl"`            // デモURL
	Frequency     string           `           json:"meetingFrequency"`   // 頻度
	StarCount     int              `           json:"starCount"`          // スター数
	Members       []Member         `           json:"currentMembers"`     // メンバー
	RecruitMember []RecruitInfo    `           json:"recruitMembers"`     // 募集メンバー
	Discussions   []DiscussionInfo `           json:"discussions"`        // 「聞いてみる」の内容
	CreateAt      time.Time        `           json:"-"`
	UpdateAt      time.Time        `           json:"-"`
}

// URLInfo 各種URL情報
type URLInfo struct {
	URL string `json:"url"`
}

type RecruitInfo struct {
	Occupation string `json:"occupation"` // 肩書とか役割
}

type Member struct {
	UserID     string `json:"id"`
	Occupation string `json:"occupation"` // 肩書とか役割
}

type DiscussionInfo struct {
	UserID    string `json:"userId"`    // ユーザー
	Message   string `json:"message"`   // コメント
	Timestamp int64  `json:"timestamp"` // 投稿日時
}

// AppsContext appsのコレクション
type AppsContext struct {
	context.Context
}

func (ctx AppsContext) Name() string {
	return "apps"
}

var _ modelsContext = (*AppsContext)(nil)

// AppsCtx appsコレクションの取得
func AppsCtx(ctx context.Context) AppsContext {
	return AppsContext{ctx}
}

func (ctx AppsContext) withCollection(fn func(c *mgo.Collection)) {
	withDefaultCollection(ctx, ctx.Name(), fn)
}

func (ctx AppsContext) FindID(appID string) (result AppInfo, err error) {
	ctx.withCollection(func(c *mgo.Collection) {
		err = c.FindId(appID).One(&result)
	})
	return
}

func (ctx AppsContext) FindAll() (result []AppInfo, err error) {
	ctx.withCollection(func(c *mgo.Collection) {
		err = c.Find(bson.M{}).All(&result)
	})
	return
}

func (ctx AppsContext) FindLatest(offset int, num int) (result []AppInfo, err error) {
	// TODO: 条件とりあえず開始日の降順（たぶん登録日にしないと）
	ctx.withCollection(func(c *mgo.Collection) {
		err = c.Find(bson.M{}).Sort("startdate").Skip(offset).Limit(num).All(&result)
	})
	return
}

func (ctx AppsContext) FindPopular(offset int, num int) (result []AppInfo, err error) {
	// TODO: 人気の条件あとで
	ctx.withCollection(func(c *mgo.Collection) {
		err = c.Find(bson.M{}).Skip(offset).Limit(num).All(&result)
	})
	return
}

// Upsert 登録
func (ctx AppsContext) Upsert(app AppInfo) error {
	var err error
	ctx.withCollection(func(c *mgo.Collection) {
		var result interface{} // bson.M
		_, err = c.FindId(app.ID).Apply(mgo.Change{
			Update: app,
			Upsert: true,
		}, &result)
	})
	return err
}

// document単位でatomicな更新
func (ctx AppsContext) findAndModify(findQuery bson.M, query bson.M) error {
	return findAndModify(ctx, findQuery, query)
}
