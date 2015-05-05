package views

import (
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
)

type AppInfoView struct {
	models.AppInfo
	Members     []UserMember      // models.Membersを上書きします
	Discussions []UserDiscussions // models.Discussionsを上書きします
	Stared      bool              // 現在認証されているユーザーがstarしているかどうか
	IsAdmin     bool              // 管理者かどうか
}

func (a AppInfoView) IsEmpty() bool {
	return a.AppInfo.ID == ""
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

	// sizeを3にする
	if len(a.ImageURLs) != 3 {
		imageURLs := make([]models.URLInfo, 3)
		for idx, img := range a.ImageURLs {
			imageURLs[idx] = img
		}
		a.ImageURLs = imageURLs
	}

	account, ok := oauth.FromContext(ctx)
	if ok {
		a.IsAdmin = a.AppInfo.IsAdmin(account.UserID)
		a.Stared = a.AppInfo.Stared(account.UserID)
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
