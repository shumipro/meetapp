package views

import (
	"net/http"

	"golang.org/x/net/context"
	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
)

/*
 - search registered user: 登録されているユーザーをkeywordマッチでリストを返す
  -> [{id: “12345”, name: “takuya tejima”},,,]
*/

func init() {
	kami.Get("/user/search/:keyword", UserSearchKeyword)
}

func UserSearchKeyword(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	word := kami.Param(ctx, "keyword")

	users := make([]models.User, 0)
	users, _ = models.UsersTable().FindByKeyword(ctx, word)

	renderer.JSON(w, 200, users)
}
