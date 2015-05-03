package views

import (
	"log"
	"net/http"

	"github.com/guregu/kami"
	"golang.org/x/net/context"
)

func init() {
	kami.Get("/u/mypage", Mypage)
}

func Mypage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := TemplateHeader{
		Title: "マイページ",
	}

	// TODO: マイページ表示時はFacebookとかの情報とりなおしてUserテーブル更新する

	if err := FromContextTemplate(ctx, "mypage").Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}
