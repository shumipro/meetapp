package constants

type ProjectState string

var projectStateMap = map[string]string{
	"1":  "募集中",
	"2":  "募集完了",
	"3":  "開発完了",
	"4":  "一時休止",
}

func init() {
	addConstants("projectState", convertConstant(projectStateMap))
}

func (t ProjectState) String() string {
	return t.Name()
}

func (t ProjectState) Name() string {
	return projectStateMap[t.ID()]
}

func (t ProjectState) ID() string {
	return string(t)
}
