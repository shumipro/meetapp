package views

import (
	"log"
	"net/http"

	"encoding/json"

	"fmt"
	"io/ioutil"

	"github.com/go-xweb/uuid"
	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
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
	kami.Get("/app/detail/:id", AppDetail)
	kami.Get("/app/list", AppList)
	kami.Get("/app/register", AppRegister)
	// API
	kami.Post("/api/app/register", AppRegisterPost)
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

type AppDetailResponse struct {
	TemplateHeader
	AppInfo models.AppInfo
}

func AppDetail(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	appID := kami.Param(ctx, "id")

	appInfo, err := models.AppsCtx(ctx).FindID(appID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}

	preload := AppListResponse{
		TemplateHeader: TemplateHeader{
			Title:    "MeetApp - " + appInfo.Name,
			SubTitle: appInfo.Name,
			NavTitle: appInfo.Name,
		},
		AppInfo: appInfo,
	}

	if err := FromContextTemplate(ctx, "app/detail").Execute(w, preload); err != nil {
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

type RegisterAppInfo struct {
	Description string  `json:"description"`
	Images      []Image `json:"images"`
	Name        string  `json:"name"`
}

type Image struct {
	URL string `json:"url"`
}

func AppRegisterPost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
	fmt.Println(string(data))

	var registerAppInfo RegisterAppInfo
	if err := json.Unmarshal(data, &registerAppInfo); err != nil {
		log.Println("ERROR! json parse", err)
		renderer.JSON(w, 400, err)
		return
	}

	// TODO: 重複チェック?

	appInfo, err := parseAppInfo(registerAppInfo)
	if err != nil {
		log.Println("ERROR! json parse", err)
		renderer.JSON(w, 400, err)
		return
	}

	if err := models.AppsCtx(ctx).Upsert(appInfo); err != nil {
		log.Println("ERROR! register", err)
		renderer.JSON(w, 400, err)
		return
	}

	renderer.JSON(w, 200, appInfo)
}

func parseAppInfo(req RegisterAppInfo) (models.AppInfo, error) {
	// TODO: 必須項目チェック?
	if req.Name == "" {
		return models.AppInfo{}, fmt.Errorf("%s", "アプリ名は必須です")
	}

	var appInfo models.AppInfo
	appInfo.ID = uuid.NewRandom().String() // TODO: とりあえずUUID
	appInfo.Name = req.Name
	appInfo.Title = req.Name // TODO: とりあえず
	appInfo.Detail = req.Description
	if len(req.Images) > 0 {
		appInfo.ImageURL = req.Images[0].URL // TODO: とりあえず1個
	}

	return appInfo, nil
}
