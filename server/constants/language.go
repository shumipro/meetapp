package constants

type LanguageType string

var pLangMap = map[string]string{
	"0":  "まだ決めていない",
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
}

func init() {
	addConstants("pLang", convertConstant(pLangMap))
}

func (t LanguageType) String() string {
	return t.Name()
}

func (t LanguageType) Name() string {
	return pLangMap[t.ID()]
}

func (t LanguageType) ID() string {
	return string(t)
}
