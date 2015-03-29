package views

import (
	"log"
	"net/http"

	"github.com/guregu/kami"
	"golang.org/x/net/context"
)

var sortLabels = map[string]map[string]string{
	"new": {
		"title": "新着アプリ",
	},
	"popular": {
		"title": "人気アプリ",
	},
}

func init() {
	kami.Get("/app/list", AppList)
	kami.Get("/app/register", AppRegister)
}

func AppList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	orderBy := r.FormValue("orderBy")
	preload := TemplateHeader{
		Title: "MeetApp - " + sortLabels[orderBy]["title"],
	}
	if err := appListTmpl.Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}

func AppRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := TemplateHeader{
		Title: "MeetApp - アプリの登録",
	}
	if err := appRegisterTmpl.Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}
