package views

import (
	"log"
	"net/http"

	"encoding/json"

	"fmt"
	"io/ioutil"

	"github.com/go-xweb/uuid"
	"github.com/guregu/kami"
	"github.com/k0kubun/pp"
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

	preload := AppListResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - "+sortLabels[orderBy]["title"],
		"",
		"気になるアプリ開発に参加しよう",
		false,
	)
	preload.AppInfoList = apps

	ExecuteTemplate(ctx, w, "app/list", preload)
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

	preload := AppDetailResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - "+appInfo.Name,
		appInfo.Name,
		appInfo.Name,
		false,
	)
	preload.AppInfo = appInfo

	ExecuteTemplate(ctx, w, "app/detail", preload)
}

func AppRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := NewHeader(ctx,
		"MeetApp - アプリの登録",
		"",
		"アプリを登録して仲間を探そう",
		false,
	)
	ExecuteTemplate(ctx, w, "app/register", preload)
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
