package constants

type AreaType string

var areaMap = map[string]string{
	"0":  "まだ決めていない",
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
}

func init() {
	addConstants("area", convertConstant(areaMap))
}

func (t AreaType) String() string {
	return t.Name()
}

func (t AreaType) Name() string {
	return areaMap[t.ID()]
}

func (t AreaType) ID() string {
	return string(t)
}