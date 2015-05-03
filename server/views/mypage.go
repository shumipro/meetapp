package views

import (
	"log"
	"net/http"

	"github.com/guregu/kami"
	"golang.org/x/net/context"
	"github.com/shumipro/meetapp/server/oauth"
	"github.com/huandu/facebook"
)

func init() {
	kami.Get("/u/mypage", Mypage)
}

func Mypage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := TemplateHeader{
		Title: "マイページ",
	}

	a, ok := oauth.FromContext(ctx)
	if !ok {
		panic("login error")
	}

	// マイページ表示時はFacebookとかの情報とりなおしてUserテーブル更新する
	_, err := facebook.Get("/me", facebook.Params{
		"access_token": a.AuthToken,
	})
	if err != nil {
		// TODO: Facebookのトークンきれたら?
		panic(err)
	}

	// TODO: 変更あればUserテーブル更新

	if err := FromContextTemplate(ctx, "mypage").Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}
