package models

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
	ID           string       // アプリID
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
