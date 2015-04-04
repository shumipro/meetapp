package models

import (
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AppCategory string

const (
	GameCategory  AppCategory = "ゲーム"
	WebCategory   AppCategory = "Webサービス"
	StudyCategory AppCategory = "学習"
)

type PlatformType string

const (
	IOS     PlatformType = "iOS"
	Android PlatformType = "Android"
	Web     PlatformType = "Web"
)

type LanguageType string

const (
	GoLang     LanguageType = "Go"
	Java       LanguageType = "Java"
	JavaScript LanguageType = "JavaScript"
	ObjectiveC LanguageType = "Objective-C"
	Swift      LanguageType = "Swift"
)

// AppInfo アプリ
type AppInfo struct {
	ID           string       `bson:"_id"` // アプリID
	Name         string       // アプリ名
	Title        string       // アプリ紹介タイトル
	Detail       string       // アプリ詳細
	URLs         string       // AppStoreのURLとかGitHubのURLとか TODO: 一旦1個あとでリストにする
	ImageURL     string       // 画像のURL（s3?）
	StarCount    int          // スター数
	Category     AppCategory  // カテゴリ
	PlatformType PlatformType // プラットフォーム
	LanguageType LanguageType // プログラミング言語
}

// URLInfo 各種URL情報
type URLInfo struct {
	Name string
	URL  string
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

func (ctx AppsContext) FindID() (result AppInfo, err error) {
	ctx.withCollection(func(c *mgo.Collection) {
		err = c.Find(bson.M{}).One(&result)
	})
	return
}

func (ctx AppsContext) FindAll() (result []AppInfo, err error) {
	ctx.withCollection(func(c *mgo.Collection) {
		err = c.Find(bson.M{}).All(&result)
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
