package constants

type MeetingFrequencyType string

var meetingFrequencyMap = map[string]string{
	"0": "まだ決めていない",
	"1": "週１回程度",
	"2": "週2, 3回程度",
	"3": "毎日",
	"4": "月1回程度",
	"5": "数ヶ月に1回程度",
	"6": "不定期",
	"7": "なし",
	"8": "その他",
}

func init() {
	addConstants("meetingFrequency", convertConstant(meetingFrequencyMap))
}

func (t MeetingFrequencyType) String() string {
	return t.Name()
}

func (t MeetingFrequencyType) Name() string {
	return meetingFrequencyMap[t.ID()]
}

func (t MeetingFrequencyType) ID() string {
	return string(t)
}
