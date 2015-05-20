package views

import (
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
	latestList, err := models.AppsInfoTable.FindLatest(ctx, 0, 4)
	if err != nil {
		panic(err)
	}
	popularList, err := models.AppsInfoTable.FindPopular(ctx, 0, 4)
	if err != nil {
		panic(err)
	}

	preload := IndexResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - 開発アイデアを実現する仲間を探そう",
		"",
		"一緒にアプリを開発する仲間を探そう",
		true,
		"",
		"",
	)
	preload.LastedList = latestList
	preload.PopularList = popularList

	ExecuteTemplate(ctx, w, r, "index", preload)
}
