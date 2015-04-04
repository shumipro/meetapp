package views

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/go-xweb/uuid"
	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
	"golang.org/x/net/context"
	"io/ioutil"
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
	kami.Post("/app/register", AppRegisterPost)
}

type AppListResponse struct {
	TemplateHeader
	AppInfoList []models.AppInfo
}

func AppList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	orderBy := r.FormValue("orderBy")

	apps, err := models.AppsCtx(ctx).FindAll()
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}

	preload := AppListResponse{
		TemplateHeader: TemplateHeader{
			Title:    "MeetApp - " + sortLabels[orderBy]["title"],
			SubTitle: "サブタイトル",
			NavTitle: "気になるアプリ開発に参加しよう",
		},
		AppInfoList: apps,
	}
	if err := FromContextTemplate(ctx, "app/list").Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}

func AppRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := TemplateHeader{
		Title:    "MeetApp - アプリの登録",
		SubTitle: "サブタイトル",
		NavTitle: "アプリを登録して仲間を探そう",
	}
	if err := FromContextTemplate(ctx, "app/register").Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}

var mockRequestData = `
{"name": "App name", "description": "hoge", "images": [{"url": "https://golang.org/doc/gopher/gopherbw.png"}]}
`

type RegisterAppInfo struct {
	Description string  `json:"description"`
	Images      []Image `json:"images"`
	Name        string  `json:"name"`
}

type Image struct {
	URL string `json:"url"`
}

func AppRegisterPost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	data := []byte(mockRequestData)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
	    return
	}

	var registerAppInfo RegisterAppInfo
	if err := json.Unmarshal(data, &registerAppInfo); err != nil {
		log.Println("ERROR! json parse", err)
		renderer.JSON(w, 400, err)
		return
	}

	// TODO: 必須項目チェック?

	// TODO: 重複チェック?

	appInfo := parseAppInfo(registerAppInfo)
	if err := models.AppsCtx(ctx).Upsert(appInfo); err != nil {
		log.Println("ERROR! register", err)
		renderer.JSON(w, 400, err)
		return
	}

	// TODO: とりあえずtop
	http.Redirect(w, r, "/", 301)
}

func parseAppInfo(req RegisterAppInfo) models.AppInfo {
	var appInfo models.AppInfo
	appInfo.ID = uuid.NewRandom().String() // TODO: とりあえずUUID
	appInfo.Name = req.Name
	appInfo.Title = req.Name // TODO: とりあえず
	appInfo.Detail = req.Description
	if len(req.Images) > 0 {
		appInfo.ImageURL = req.Images[0].URL // TODO: とりあえず1個
	}
	return appInfo
}
