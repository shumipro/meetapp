package views

import (
	"net/http"

	"log"
	"strings"

	"github.com/guregu/kami"
	"github.com/huandu/facebook"
	"github.com/kyokomi/cloudinary"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
)

func init() {
	kami.Get("/mypage/other/:id", MypageOther)

	kami.Get("/u/mypage", Mypage)
	// API
	kami.Get("/u/api/upload/image", UploadImage)
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

	ExecuteTemplate(ctx, w, "mypage", preload)
}

func UploadImage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := oauth.FromContext(ctx)

	formFile, _, err := r.FormFile("file")
	if err != nil {
		renderer.JSON(w, 400, err)
		return
	}
	defer formFile.Close()

	if err := cloudinary.UploadStaticImage(ctx, a.UserID, formFile); err != nil {
		renderer.JSON(w, 400, err)
		return
	}

	largeImageURL := cloudinary.ResourceURL(ctx, a.UserID)
	user, err := models.UsersTable.FindID(ctx, a.UserID)
	if err != nil {
		renderer.JSON(w, 400, err)
		return
	}

	user.LargeImageURL = largeImageURL
	user.ImageURL = strings.Replace(largeImageURL, "image/upload", "image/upload/w_96,h_96", 1)

	if err := models.UsersTable.Upsert(ctx, user); err != nil {
		renderer.JSON(w, 400, err)
		return
	}

	renderer.JSON(w, 200, largeImageURL)
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

	ExecuteTemplate(ctx, w, "mypage", preload)
}
