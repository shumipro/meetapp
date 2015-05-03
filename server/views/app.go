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
	"github.com/k0kubun/pp"
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
	kami.Post("/api/app/discussion", AppDiscussionPost)
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

	pp.Println(registerAppInfo)

	// TODO: 重複チェック?

	registerAppInfo.ID = uuid.NewRandom().String()
	if len(registerAppInfo.ImageURLs) > 0 {
		registerAppInfo.MainImage = registerAppInfo.ImageURLs[0].URL // TODO: とりあえず1件目をメインの画像にする
	} else {
		// set default image
		registerAppInfo.MainImage = "/img/no_img.png"
	}

	if err := models.AppsCtx(ctx).Upsert(registerAppInfo); err != nil {
		log.Println("ERROR! register", err)
		renderer.JSON(w, 400, err)
		return
	}

	renderer.JSON(w, 200, registerAppInfo)
}

type DiscussionRequest struct {
	AppID string `json:"appId"`     // アプリID
	DiscussionInfo models.DiscussionInfo
}

func AppDiscussionPost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
	fmt.Println(string(data))

	// convert request params to struct
	var discussionReq DiscussionRequest
	if err := json.Unmarshal(data, &discussionReq); err != nil {
		log.Println("ERROR! json parse", err)
		renderer.JSON(w, 400, err)
		return
	}

	// get appinfo from db
	appInfo, err := models.AppsCtx(ctx).FindID(discussionReq.AppID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}

	// push a discussionInfo
	appInfo.Discussions = append(appInfo.Discussions, discussionReq.DiscussionInfo) 

	if err := models.AppsCtx(ctx).Upsert(appInfo); err != nil {
		log.Println("ERROR! register", err)
		renderer.JSON(w, 400, err)
		return
	}

	renderer.JSON(w, 200, appInfo.Discussions)
}
