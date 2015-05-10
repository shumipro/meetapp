package models

import "gopkg.in/mgo.v2/bson"

type AppInfoOrderType string

const (
	OrderByNew      AppInfoOrderType = "new"
	OrderByPopular  AppInfoOrderType = "popular"
	OrderByUpdateAt AppInfoOrderType = "updateAt"
)

type AppInfoFilter struct {
	OccupationType OccupationType
	CategoryType   CategoryType
	PlatformType   PlatformType
	LanguageType   LanguageType
	AreaType       AreaType
	OrderBy        AppInfoOrderType
}

func (a AppInfoOrderType) SortCondition() string {
	if a == OrderByUpdateAt {
		return "-updateat"
	} else if a == OrderByPopular {
		return "-starcount"
	} else if a == OrderByNew {
		return "-createat"
	}
	return ""
}

func (a AppInfoFilter) Condition() bson.M {
	condition := bson.M{}

	// AppInfo
	if a.CategoryType != "" {
		condition["category"] = a.CategoryType
	}
	if a.PlatformType != "" {
		condition["platform"] = a.PlatformType
	}
	if a.LanguageType != "" {
		condition["plang"] = a.LanguageType
	}
	if a.AreaType != "" {
		condition["area"] = a.AreaType
	}

	// AppInfo.RecruitMember
	if a.OccupationType != "" {
		condition["recruitmember.occupation"] = a.OccupationType
	}

	return condition
}
