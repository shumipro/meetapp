package views

import (
	"net/http"

	"log"

	"errors"
	"time"

	"github.com/guregu/kami"
	"github.com/kyokomi/cloudinary"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
)

/*
 - search registered user: 登録されているユーザーをkeywordマッチでリストを返す
  -> [{id: “12345”, name: “takuya tejima”},,,]
*/

func init() {
	kami.Put("/u/api/user", UserProfileUpdate)
	kami.Get("/api/user/search/:keyword", UserSearchKeyword)
}

func UserSearchKeyword(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	word := kami.Param(ctx, "keyword")

	users, _ := models.UsersTable.FindByKeyword(ctx, word)
	if users == nil {
		users = []models.User{}
	}

	renderer.JSON(w, 200, users)
}

// UserProfileUpdate ユーザーのプロフィール更新API
func UserProfileUpdate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := oauth.FromContext(ctx)

	user, err := readBodyUser(r.Body)
	if err != nil {
		renderer.JSON(w, 400, err.Error())
		return
	}

	if a.UserID != user.ID {
		renderer.JSON(w, 401, errors.New("loginしたユーザと一致しません").Error())
		return
	}

	beforeUser, err := models.UsersTable.FindID(ctx, a.UserID)
	if err != nil {
		renderer.JSON(w, 400, err.Error())
		return
	}

	// 前の画像を削除する
	if beforeUser.ImageName != "" {
		if err := cloudinary.DeleteStaticImage(ctx, beforeUser.ImageName); err != nil {
			// 失敗時にログだけ出す
			log.Println(err)
		}
	}

	// Userテーブル更新する（一部は以前の情報をそのまま）
	user.FBUser = beforeUser.FBUser
	user.CreateAt = beforeUser.CreateAt
	user.UpdateAt = time.Now()
	if err := models.UsersTable.Upsert(ctx, user); err != nil {
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, user)
}
