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

	preload := AppDetailResponse{
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

func AppRegisterPost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
	fmt.Println(string(data))

	var registerAppInfo models.AppInfo
	if err := json.Unmarshal(data, &registerAppInfo); err != nil {
		log.Println("ERROR! json parse", err)
		renderer.JSON(w, 400, err)
		return
	}

	// TODO: 重複チェック?

	registerAppInfo.ID = uuid.NewRandom().String()
	if len(registerAppInfo.ImageURLs) > 0 {
		registerAppInfo.MainImage = registerAppInfo.ImageURLs[0].URL // TODO: とりあえず1件目をメインの画像にする
	} else {
		// set default image
		registerAppInfo.MainImage = "/img/no_img.png"
	}
	// TODO: memberと募集はrequestにないので一旦固定値
	registerAppInfo.Members = []models.Member{
		{
			Name:         "kyokomi",
			IconImageURL: "https://avatars0.githubusercontent.com/u/1456047?v=3&s=460",
			Post:         "Gopher",
		},
		{
			Name:         "tejitak",
			IconImageURL: "http://graph.facebook.com/10152160532855662/picture?type=square",
			Post:         "Engineer",
		},
	}
	registerAppInfo.RecruitMember = []models.RecruitInfo{
		{
			Post: "デザイナー",
			Num:  1,
		},
		{
			Post: "企画",
			Num:  1,
		},
	}

	if err := models.AppsCtx(ctx).Upsert(registerAppInfo); err != nil {
		log.Println("ERROR! register", err)
		renderer.JSON(w, 400, err)
		return
	}

	renderer.JSON(w, 200, registerAppInfo)
}
