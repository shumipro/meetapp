package models

// AppInfo アプリ
type AppInfo struct {
	ID        string // アプリID
	Name      string // アプリ名
	Title     string // アプリ紹介タイトル
	Detail    string // アプリ詳細
	URLs      string // AppStoreのURLとかGitHubのURLとか TODO: 一旦1個あとでリストにする
	ImageURL  string // 画像のURL（s3?）
	StarCount int    // スター数
}

// URLInfo 各種URL情報
type URLInfo struct {
	Name string
	URL  string
}
