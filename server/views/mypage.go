package views

import (
	"net/http"

	"github.com/guregu/kami"
	"github.com/huandu/facebook"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
	"github.com/shumipro/meetapp/server/models"
)

func init() {
	kami.Get("/u/mypage", Mypage)
}

type MyPageResponse struct {
	TemplateHeader
	AdminAppList []AppInfoView
	JoinAppList  []AppInfoView
}

func Mypage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, ok := oauth.FromContext(ctx)
	if !ok {
		oauth.ResetCacheAuthToken(ctx, w)
		panic("login error")
	}

	// マイページ表示時はFacebookとかの情報とりなおしてUserテーブル更新する
	_, err := facebook.Get("/me", facebook.Params{
		"access_token": a.AuthToken,
	})
	if err != nil {
		oauth.ResetCacheAuthToken(ctx, w)
		panic(err)
	}

	// TODO: Facebook情報に変更あればUserテーブル更新する

	preload := MyPageResponse{}
	preload.TemplateHeader = NewHeader(ctx, "マイページ", "", "", false)

	adminApps, _ := models.AppsInfoTable.FindByAdminID(ctx, a.UserID)
	joinApps, _ := models.AppsInfoTable.FindByJoinID(ctx, a.UserID)
	preload.AdminAppList = convertAppInfoViewList(ctx, adminApps)
	preload.JoinAppList = convertAppInfoViewList(ctx, joinApps)

	ExecuteTemplate(ctx, w, "mypage", preload)
}
