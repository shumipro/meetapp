package views

import (
	"net/http"

	"github.com/guregu/kami"
	"github.com/huandu/facebook"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
)

func init() {
	kami.Get("/u/mypage", Mypage)
}

func Mypage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
	preload := NewHeader(ctx, "マイページ", "", "", false)
	ExecuteTemplate(ctx, w, "mypage", preload)
}
