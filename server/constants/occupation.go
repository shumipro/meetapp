package constants

type OccupationType string

var occupationMap = map[string]string{
	"0":  "まだ決めていない",
	"1":  "企画",
	"2":  "デザイナー",
	"3":  "サーバーエンジニア",
	"4":  "スマートフォンエンジニア",
	"5":  "フロントエンドエンジニア",
	"6":  "Webエンジニア",
	"7":  "インフラエンジニア",
	"8":  "フルスタックエンジニア",
	"9":  "プロジェクトマネージャー",
	"10": "営業",
	"11": "おてつだい",
	"12": "その他",
}

func init() {
	addConstants("occupation", convertConstant(occupationMap))
}

func (t OccupationType) String() string {
	return t.Name()
}

func (t OccupationType) Name() string {
	return occupationMap[t.ID()]
}

func (t OccupationType) ID() string {
	return string(t)
}
