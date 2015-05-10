package views

import (
	"log"
	"net/http"

	"github.com/guregu/kami"
	"github.com/huandu/facebook"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
)

func init() {
	kami.Get("/mypage/other/:id", MypageOther)

	kami.Get("/u/mypage", Mypage)
	kami.Get("/u/mypage/edit", MypageEdit)
}

type MyPageResponse struct {
	TemplateHeader
	User         models.User
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

	user, err := models.UsersTable.FindID(ctx, a.UserID)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	preload := MyPageResponse{}
	preload.User = user
	preload.TemplateHeader = NewHeader(ctx, "マイページ", "", "", false)

	adminApps, _ := models.AppsInfoTable.FindByAdminID(ctx, a.UserID)
	joinApps, _ := models.AppsInfoTable.FindByJoinID(ctx, a.UserID)
	preload.AdminAppList = convertAppInfoViewList(ctx, adminApps)
	preload.JoinAppList = convertAppInfoViewList(ctx, joinApps)

	ExecuteTemplate(ctx, w, r, "mypage", preload)
}

func MypageOther(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID := kami.Param(ctx, "id")
	if userID == "" || userID == "favicon.png" {
		return
	}

	user, err := models.UsersTable.FindID(ctx, userID)
	if err != nil {
		log.Println(err, userID)
		panic(err)
	}

	preload := MyPageResponse{}
	preload.User = user
	preload.TemplateHeader = NewHeader(ctx, user.Name, "", "", false)

	joinApps, _ := models.AppsInfoTable.FindByJoinID(ctx, userID)
	preload.JoinAppList = convertAppInfoViewList(ctx, joinApps)

	ExecuteTemplate(ctx, w, r, "mypage", preload)
}

func MypageEdit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, ok := oauth.FromContext(ctx)
	if !ok {
		oauth.ResetCacheAuthToken(ctx, w)
		panic("login error")
	}

	user, err := models.UsersTable.FindID(ctx, a.UserID)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	preload := MyPageResponse{}
	preload.User = user
	preload.TemplateHeader = NewHeader(ctx, "マイページの編集", "", "", false)

	ExecuteTemplate(ctx, w, r, "mypageEdit", preload)
}
