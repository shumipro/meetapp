package constants

type CategoryType string

var categoryMap = map[string]string{
	"0":  "まだ決めていない",
	"1":  "ブック",
	"2":  "ビジネス",
	"3":  "カタログ",
	"4":  "教育",
	"5":  "エンターテイメント",
	"6":  "ファイナンス",
	"7":  "フード",
	"8":  "ゲーム",
	"9":  "ヘルスケア",
	"10": "ライフスタイル",
	"11": "メディカル",
	"12": "ミュージック",
	"13": "ナビゲーション",
	"14": "ニュース",
	"15": "写真/ビデオ",
	"16": "仕事効率化",
	"17": "辞書/辞典",
	"18": "SNS",
	"19": "スポーツ",
	"20": "旅行",
	"21": "ユーティリティ",
	"22": "天気",
	"23": "その他",
}

func init() {
	addConstants("category", convertConstant(categoryMap))
}

func (t CategoryType) String() string {
	return t.Name()
}

func (t CategoryType) Name() string {
	return categoryMap[t.ID()]
}

func (t CategoryType) ID() string {
	return string(t)
}
