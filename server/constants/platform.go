package constants

type PlatformType string

var platformMap = map[string]string{
	"0": "まだ決めていない",
	"1": "Web",
	"2": "iOS",
	"3": "Android",
	"4": "Mac",
	"5": "Apple Watch",
	"6": "Windows",
	"7": "Linux",
	"8": "Chrome",
	"9": "その他",
}

func init() {
	addConstants("platform", convertConstant(platformMap))
}

func (t PlatformType) String() string {
	return t.Name()
}

func (t PlatformType) Name() string {
	return platformMap[t.ID()]
}

func (t PlatformType) ID() string {
	return string(t)
}