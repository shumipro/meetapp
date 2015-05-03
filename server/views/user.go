package views

import (
	"net/http"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
	"golang.org/x/net/context"
)

/*
 - search registered user: 登録されているユーザーをkeywordマッチでリストを返す
  -> [{id: “12345”, name: “takuya tejima”},,,]
*/

func init() {
	kami.Get("/api/user/search/:keyword", UserSearchKeyword)
}

func UserSearchKeyword(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	word := kami.Param(ctx, "keyword")

	users, _ := models.UsersTable().FindByKeyword(ctx, word)
	if users == nil {
		users = []models.User{}
	}

	renderer.JSON(w, 200, users)
}
