package models

import (
	"sort"

	"github.com/mattn/natural"
)

type Constants map[string][]Constant

type Constant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ConstantSlice []Constant

func (p ConstantSlice) Len() int           { return len(p) }
func (p ConstantSlice) Less(i, j int) bool { return natural.NaturalCaseComp(p[i].ID, p[j].ID) < 0 }
func (p ConstantSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ConstantSlice) Sort()              { sort.Sort(p) }

func AllConstants() Constants {
	c := map[string][]Constant{}
	c["platform"] = convetConstant(platformMap)
	c["category"] = convetConstant(categoryMap)
	c["pLang"] = convetConstant(pLangMap)
	c["area"] = convetConstant(areaMap)
	c["occupation"] = convetConstant(occupationMap)
	c["meetingFrequency"] = convetConstant(meetingFrequencyMap)
	return Constants(c)
}

func convetConstant(m map[string]string) []Constant {
	c := []Constant{}
	for key, val := range m {
		c = append(c, Constant{key, val})
	}

	sort.Sort(ConstantSlice(c))
	return c
}

type PlatformType string

func (t PlatformType) Name() string {
	return platformMap[string(t)]
}

var platformMap = map[string]string{
	"1":  "Web",
	"2":  "iOS",
	"3":  "Android",
	"4":  "Mac",
	"5":  "Apple Watch",
	"6":  "Windows",
	"7":  "Linux",
	"8":  "Chrome",
	"9":  "その他",
	"99": "まだ決めていない",
}

type CategoryType string

var categoryMap = map[string]string{
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
	"99": "まだ決めていない",
}

func (t CategoryType) Name() string {
	return categoryMap[string(t)]
}

type LanguageType string

var pLangMap = map[string]string{
	"1":  "Java",
	"2":  "JavaScript",
	"3":  "Swift",
	"4":  "Objective-c",
	"5":  "Go",
	"6":  "PHP",
	"7":  "Python",
	"8":  "Ruby",
	"9":  "Perl",
	"10": "C",
	"11": "C#",
	"12": "C++",
	"13": "Scala",
	"14": "SQL",
	"15": "VB",
	"16": "MATLAB",
	"17": "その他",
	"99": "まだ決めていない",
}

func (t LanguageType) Name() string {
	return pLangMap[string(t)]
}

type AreaType string

var areaMap = map[string]string{
	"1":  "渋谷",
	"2":  "六本木",
	"3":  "恵比寿",
	"4":  "新宿",
	"5":  "横浜",
	"6":  "池袋",
	"7":  "品川",
	"8":  "五反田",
	"9":  "有楽町",
	"10": "大阪",
	"11": "福岡",
	"12": "その他",
	"99": "まだ決めていない",
}

func (t AreaType) Name() string {
	return areaMap[string(t)]
}

type OccupationType string

var occupationMap = map[string]string{
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
	"99": "まだ決めていない",
}

func (t OccupationType) Name() string {
	return occupationMap[string(t)]
}

type MeetingFrequencyType string

var meetingFrequencyMap = map[string]string{
	"1":  "週１回程度",
	"2":  "週2, 3回程度",
	"3":  "毎日",
	"4":  "月1回程度",
	"5":  "数ヶ月に1回程度",
	"6":  "不定期",
	"7":  "なし",
	"8":  "その他",
	"99": "まだ決めていない",
}

func (t MeetingFrequencyType) Name() string {
	return meetingFrequencyMap[string(t)]
}
