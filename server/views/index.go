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

	preload := IndexResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp",
		"サブタイトル",
		"一緒にアプリを開発する仲間を探そう",
		true,
	)
	preload.LastedList = latestList
	preload.PopularList = popularList

	ExecuteTemplate(ctx, w, "index", preload)
}
