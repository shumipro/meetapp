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
}

type IndexResponse struct {
	TemplateHeader
	LastedList  []models.AppInfo // 新着アプリ
	PopularList []models.AppInfo // 人気アプリ
}

var mockDataList = []models.AppInfo{
	{
		ID:       "001",
		Name:     "アプリ1",
		Title:    "1あああああああああ",
		Detail:   "1いいいいいいいいいいいいいいいいいい\nいいいいいいい",
		URLs:     "http://google.com",
		ImageURL: "https://golang.org/doc/gopher/gopherbw.png",
		StarCount: 100,
	},
	{
		ID:       "002",
		Name:     "アプリ2",
		Title:    "2あああああああああ",
		Detail:   "2いいいいいいいいいいいいいいいいいい\nいいいいいいい",
		URLs:     "http://google.com",
		ImageURL: "https://golang.org/doc/gopher/gopherbw.png",
		StarCount: 30,
	},
	{
		ID:       "003",
		Name:     "アプリ3",
		Title:    "3あああああああああ",
		Detail:   "3いいいいいいいいいいいいいいいいいい\nいいいいいいい",
		URLs:     "http://google.com",
		ImageURL: "https://golang.org/doc/gopher/gopherbw.png",
		StarCount: 20,
	},
	{
		ID:       "004",
		Name:     "アプリ4",
		Title:    "4あああああああああ",
		Detail:   "4いいいいいいいいいいいいいいいいいい\nいいいいいいい",
		URLs:     "http://google.com",
		ImageURL: "https://golang.org/doc/gopher/gopherbw.png",
		StarCount: 40,
	},
}

func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := IndexResponse{
		TemplateHeader: TemplateHeader{
			Title: "MeetApp",
			SubTitle: "サブタイトル",
			ShowBanner: true,
		},
		LastedList:  mockDataList,
		PopularList: mockDataList,
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
