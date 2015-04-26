package views

import (
	"log"
	"net/http"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
	"golang.org/x/net/context"
)

func init() {
	kami.Get("/", Index)
	kami.Get("/error", Error)
	kami.Get("/about", About)
	kami.Get("/login", Login)
	kami.Get("/mypage", Mypage)
}

type IndexResponse struct {
	TemplateHeader
	LastedList  []models.AppInfo // 新着アプリ
	PopularList []models.AppInfo // 人気アプリ
}

func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	latestList, err := models.AppsCtx(ctx).FindLatest(0, 4)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
	popularList, err := models.AppsCtx(ctx).FindPopular(0, 4)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}

	preload := IndexResponse{
		TemplateHeader: TemplateHeader{
			Title:      "MeetApp",
			SubTitle:   "サブタイトル",
			NavTitle:   "一緒にアプリを開発する仲間を探そう",
			ShowBanner: true,
		},
		LastedList:  latestList,
		PopularList: popularList,
	}
	if err := FromContextTemplate(ctx, "index").Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}

func Error(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := TemplateHeader{
		Title: "Error",
	}
	if err := FromContextTemplate(ctx, "error").Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}

func About(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := TemplateHeader{
		Title: "About",
	}
	if err := FromContextTemplate(ctx, "about").Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}

func Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := TemplateHeader{
		Title: "Login",
	}
	if err := FromContextTemplate(ctx, "login").Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}

func Mypage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := TemplateHeader{
		Title: "マイページ",
	}
	if err := FromContextTemplate(ctx, "mypage").Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}
