package views

import (
	"github.com/shumipro/meetapp/server/models"
	"golang.org/x/net/context"
)

type AppInfoView struct {
	models.AppInfo
	Members     []UserMember      // models.Membersを上書きします
	Discussions []UserDiscussions // models.Discussionsを上書きします
}

// UserMember User情報を持つMember
type UserMember struct {
	models.Member
	models.User
}

type UserDiscussions struct {
	models.DiscussionInfo
	models.User
}

func NewAppInfoView(ctx context.Context, appInfo models.AppInfo) AppInfoView {
	a := AppInfoView{}
	a.AppInfo = appInfo

	a.Members = make([]UserMember, len(a.AppInfo.Members))
	for idx, m := range appInfo.Members {
		// TODO: あとでIn句にして1クエリにする
		u, _ := models.UsersTable.FindID(ctx, m.UserID)
		a.Members[idx] = UserMember{Member: m, User: u}
	}

	a.Discussions = make([]UserDiscussions, len(a.AppInfo.Discussions))
	for idx, d := range appInfo.Discussions {
		// TODO: あとでIn句にして1クエリにする
		u, _ := models.UsersTable.FindID(ctx, d.UserID)
		a.Discussions[idx] = UserDiscussions{DiscussionInfo: d, User: u}
	}
	return a
}

func convertAppInfoViewList(ctx context.Context, apps []models.AppInfo) []AppInfoView {
	appViews := make([]AppInfoView, len(apps))
	for idx, app := range apps {
		appViews[idx] = NewAppInfoView(ctx, app)
	}
	return appViews
}
