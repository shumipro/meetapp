package views

import (
	"log"
	"net/http"

	"encoding/json"

	"fmt"
	"io/ioutil"

	"strings"
	"time"

	"strconv"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
)

var notAdminError = fmt.Errorf("%s", "not admin user")

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
	kami.Get("/u/app/register", AppRegister)
	kami.Get("/u/app/edit/:id", AppEdit)
	// Apps API
	//	kami.Get("/u/api/app/apps", APIAppGetAll)
	//	kami.Get("/u/api/app/apps/:id", APIAppGet)
	kami.Post("/u/api/app/register", APIAppRegister)   // TODO: [POST] /u/api/app/apps
	kami.Put("/u/api/app/edit/:id", APIAppEdit)        // TODO: [PUT] /u/api/app/apps/:id
	kami.Delete("/u/api/app/delete/:id", APIAppDelete) // TODO: [DELETE] /u/api/app/apps/:id
	// Discussion API
	kami.Post("/u/api/app/discussion", APIAppDiscussion)
	// Star API
	kami.Post("/u/api/app/star/:id", APIAppStared)
	kami.Delete("/u/api/app/star/:id", APIAppStarDelete)
}

type AppListResponse struct {
	TemplateHeader
	AppInfoList []AppInfoView
	CurrentPage int
	PerPageNum  int
	TotalCount  int
}

const perPageNum int = 10

func AppList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	orderBy := r.FormValue("orderBy")

	page, _ := strconv.Atoi(r.FormValue("page"))
	if page > 0 {
		page -= 1 // 1ページは0とする
	}

	platform := r.FormValue("platform")
	occupation := r.FormValue("occupation")
	category := r.FormValue("category")
	pLang := r.FormValue("pLang")
	area := r.FormValue("area")

	filter := models.AppInfoFilter{}
	filter.OccupationType = models.OccupationType(occupation)
	filter.CategoryType = models.CategoryType(category)
	filter.LanguageType = models.LanguageType(pLang)
	filter.AreaType = models.AreaType(area)
	filter.PlatformType = models.PlatformType(platform)

	preload := AppListResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - "+sortLabels[orderBy]["title"],
		"",
		"気になるアプリ開発に参加しよう",
		false,
	)

	// ViewModel変換して詰める
	totalCount, apps, err := models.AppsInfoTable.FindFilter(ctx, filter, page*perPageNum, perPageNum)
	if err != nil {
		panic(err)
	}
	preload.AppInfoList = convertAppInfoViewList(ctx, apps)
	preload.PerPageNum = perPageNum
	preload.CurrentPage = page + 1
	preload.TotalCount = totalCount

	ExecuteTemplate(ctx, w, "app/list", preload)
}

type AppDetailResponse struct {
	TemplateHeader
	AppInfo AppInfoView
}

func AppDetail(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	appID := kami.Param(ctx, "id")
	// TODO: とりあえず
	if appID == "favicon.png" || appID == "" {
		return
	}
	appInfo, err := models.AppsInfoTable.FindID(ctx, appID)
	if err != nil {
		panic(err)
	}

	preload := AppDetailResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - "+appInfo.Name,
		appInfo.Name,
		appInfo.Name,
		false,
	)
	preload.AppInfo = NewAppInfoView(ctx, appInfo)

	ExecuteTemplate(ctx, w, "app/detail", preload)
}

type AppRegisterResponse struct {
	TemplateHeader
	AppInfo AppInfoView
}

func AppRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := AppRegisterResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - アプリの登録",
		"",
		"アプリを登録して仲間を探そう",
		false,
	)

	// 自分をデフォルトメンバーとして突っ込んでおく
	a, _ := oauth.FromContext(ctx)
	appInfo := models.AppInfo{}
	members := []models.Member{
		{
			UserID:     a.UserID,
			Occupation: models.OccupationType("1"),
			IsAdmin:    true,
		},
	}
	appInfo.Members = members
	// sizeを3にする
	appInfo.ImageURLs = make([]models.URLInfo, 3)

	preload.AppInfo = NewAppInfoView(ctx, appInfo)

	ExecuteTemplate(ctx, w, "app/register", preload)
}

func AppEdit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	appID := kami.Param(ctx, "id")

	// TODO: とりあえず
	if appID == "favicon.png" || appID == "" {
		return
	}

	appInfo, err := models.AppsInfoTable.FindID(ctx, appID)
	if err != nil {
		panic(err)
	}

	preload := AppRegisterResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - アプリの編集",
		"",
		"アプリを登録して仲間を探そう",
		false,
	)
	// sizeを3にする
	if len(appInfo.ImageURLs) != 3 {
		imageURLs := make([]models.URLInfo, 3)
		for idx, img := range appInfo.ImageURLs {
			imageURLs[idx] = img
		}
		appInfo.ImageURLs = imageURLs
	}
	preload.AppInfo = NewAppInfoView(ctx, appInfo)

	ExecuteTemplate(ctx, w, "app/register", preload)
}

func APIAppRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	regAppInfo, err := readBodyAppInfo(r.Body)
	if err != nil {
		renderer.JSON(w, 400, "[ERROR] request param appInfo "+err.Error())
		return
	}

	regAppInfo = convertRegisterAppInfo(ctx, regAppInfo)

	if err := models.AppsInfoTable.Upsert(ctx, regAppInfo); err != nil {
		log.Println("ERROR! register", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, regAppInfo)
}

func APIAppEdit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	regAppInfo, err := readBodyAppInfo(r.Body)
	if err != nil || regAppInfo.ID == "" {
		renderer.JSON(w, 400, "[ERROR] request param appInfo ")
		return
	}


	// アプリが存在しないか
	beforeApp, err := models.AppsInfoTable.FindID(ctx, regAppInfo.ID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// 管理者じゃないアプリか
	a, _ := oauth.FromContext(ctx)
	if !beforeApp.IsAdmin(a.UserID) {
		log.Println("ERROR!", notAdminError)
		renderer.JSON(w, 400, notAdminError.Error())
		return
	}

	regAppInfo = convertEditAppInfo(ctx, regAppInfo, beforeApp)

	if err := models.AppsInfoTable.Upsert(ctx, regAppInfo); err != nil {
		log.Println("ERROR! register", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, regAppInfo)
}

func APIAppDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := oauth.FromContext(ctx)
	appID := kami.Param(ctx, "id")

	app, err := models.AppsInfoTable.FindID(ctx, appID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// 管理者のみ削除可能
	if !app.IsAdmin(a.UserID) {
		log.Println("ERROR!", notAdminError)
		renderer.JSON(w, 400, notAdminError.Error())
		return
	}

	if err := models.AppsInfoTable.Delete(ctx, app.ID); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, appID)
}

type DiscussionRequest struct {
	AppID          string `json:"appId"` // アプリID
	DiscussionInfo models.DiscussionInfo
}

func APIAppDiscussion(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}
	fmt.Println(string(data))

	// convert request params to struct
	var discussionReq DiscussionRequest
	if err := json.Unmarshal(data, &discussionReq); err != nil {
		log.Println("ERROR! json parse", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// get appinfo from db
	appInfo, err := models.AppsInfoTable.FindID(ctx, discussionReq.AppID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	nowTime := time.Now()
	discussionReq.DiscussionInfo.Timestamp = nowTime
	// push a discussionInfo
	appInfo.Discussions = append(appInfo.Discussions, discussionReq.DiscussionInfo)
	appInfo.UpdateAt = nowTime

	if err := models.AppsInfoTable.Upsert(ctx, appInfo); err != nil {
		log.Println("ERROR! discussion", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, appInfo.Discussions)
}

func APIAppStared(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := oauth.FromContext(ctx)
	appID := kami.Param(ctx, "id")

	// get appinfo from db
	appInfo, err := models.AppsInfoTable.FindID(ctx, appID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// すでにスター済み
	if appInfo.Stared(a.UserID) {
		log.Println("WARN", "stared")
		renderer.JSON(w, 200, appInfo.StarUsers)
		return
	}

	// push the user as starUsers
	appInfo.StarUsers = append(appInfo.StarUsers, a.UserID)
	// update starCount
	appInfo.StarCount = len(appInfo.StarUsers)

	appInfo.UpdateAt = time.Now()

	if err := models.AppsInfoTable.Upsert(ctx, appInfo); err != nil {
		log.Println("ERROR! star", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, appInfo.StarUsers)
}

func APIAppStarDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := oauth.FromContext(ctx)
	appID := kami.Param(ctx, "id")

	// get appinfo from db
	appInfo, err := models.AppsInfoTable.FindID(ctx, appID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// すでに削除済み
	if !appInfo.Stared(a.UserID) {
		log.Println("WARN!", "not stared")
		renderer.JSON(w, 200, appInfo.StarUsers)
		return
	}

	for idx, userID := range appInfo.StarUsers {
		if !strings.EqualFold(userID, userID) {
			continue
		}
		// remove the user from starUsers list
		appInfo.StarUsers = append(appInfo.StarUsers[:idx], appInfo.StarUsers[idx+1:]...)
		// update starCount
		appInfo.StarCount = len(appInfo.StarUsers)
		break
	}

	appInfo.UpdateAt = time.Now()

	if err := models.AppsInfoTable.Upsert(ctx, appInfo); err != nil {
		log.Println("ERROR! star", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, appInfo.StarUsers)
}
