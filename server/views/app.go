package views

import (
	"log"
	"net/http"

	"encoding/json"

	"fmt"
	"io/ioutil"
	
	"strings"
	"time"

	"github.com/go-xweb/uuid"
	"github.com/guregu/kami"
	"github.com/k0kubun/pp"
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
	// API
	kami.Post("/u/api/app/register", AppRegisterPost)
	kami.Post("/u/api/app/discussion", AppDiscussionPost)
	kami.Delete("/u/api/app/delete/:id", AppDelete)
	kami.Post("/u/api/app/star", AppStarPost)
	kami.Delete("/u/api/app/star", AppStarDelete)
}

type AppListResponse struct {
	TemplateHeader
	AppInfoList []AppInfoView
}

func AppList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	orderBy := r.FormValue("orderBy")

	preload := AppListResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - "+sortLabels[orderBy]["title"],
		"",
		"気になるアプリ開発に参加しよう",
		false,
	)

	// ViewModel変換して詰める
	apps, err := models.AppsInfoTable.FindAll(ctx)
	if err != nil {
		panic(err)
	}
	preload.AppInfoList = convertAppInfoViewList(ctx, apps)

	ExecuteTemplate(ctx, w, "app/list", preload)
}

type AppDetailResponse struct {
	TemplateHeader
	AppInfo AppInfoView
}

func (a AppDetailResponse) IsAdmin() bool {
	if a.Config.User.IsEmpty() {
		return false
	}
	return a.AppInfo.IsAdmin(a.Config.User.ID)
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

func AppRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := NewHeader(ctx,
		"MeetApp - アプリの登録",
		"",
		"アプリを登録して仲間を探そう",
		false,
	)
	ExecuteTemplate(ctx, w, "app/register", preload)
}

type AppEditResponse struct {
	TemplateHeader
	AppInfo AppInfoView
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

	preload := AppEditResponse{}
	preload.TemplateHeader = NewHeader(ctx,
		"MeetApp - アプリの編集",
		"",
		"アプリを登録して仲間を探そう",
		false,
	)
	preload.AppInfo = NewAppInfoView(ctx, appInfo)

	ExecuteTemplate(ctx, w, "app/register", preload)
}

func AppRegisterPost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}
	fmt.Println(string(data))

	var regAppInfo models.AppInfo
	if err := json.Unmarshal(data, &regAppInfo); err != nil {
		log.Println("ERROR! json parse", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	a, _ := oauth.FromContext(ctx)

	// 登録時刻、更新時刻
	nowTime := time.Now()

	// すでに登録済み
	if regAppInfo.ID != "" {
		app, err := models.AppsInfoTable.FindID(ctx, regAppInfo.ID)
		if err != nil {
			log.Println("ERROR!", err)
			renderer.JSON(w, 400, err.Error())
			return
		}

		// 管理者じゃない
		if !app.IsAdmin(a.UserID) {
			log.Println("ERROR!", notAdminError)
			renderer.JSON(w, 400, notAdminError.Error())
			return
		}
		// 登録日だけ残して後は上書きする
		regAppInfo.CreateAt = app.CreateAt
	} else {
		regAppInfo.CreateAt = nowTime
	}

	regAppInfo.UpdateAt = nowTime

	// 管理者設定
	for idx, m := range regAppInfo.Members {
		if m.UserID != a.UserID {
			continue
		}
		regAppInfo.Members[idx].IsAdmin = true
	}

	// メインの画像を設定
	regAppInfo.ID = uuid.NewRandom().String()
	if len(regAppInfo.ImageURLs) > 0 {
		regAppInfo.MainImage = regAppInfo.ImageURLs[0].URL // TODO: とりあえず1件目をメインの画像にする
	} else {
		// set default image
		regAppInfo.MainImage = "/img/no_img.png"
	}

	pp.Println(regAppInfo)

	if err := models.AppsInfoTable.Upsert(ctx, regAppInfo); err != nil {
		log.Println("ERROR! register", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, regAppInfo)
}

func AppDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

func AppDiscussionPost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	// push a discussionInfo
	appInfo.Discussions = append(appInfo.Discussions, discussionReq.DiscussionInfo)
	nowTime := time.Now()
	appInfo.UpdateAt = nowTime

	if err := models.AppsInfoTable.Upsert(ctx, appInfo); err != nil {
		log.Println("ERROR! discussion", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, appInfo.Discussions)
}

func indexOf(s []string, e string) int {
    for i, a := range s { if strings.EqualFold(a, e) { return i } }
    return -1
}

type StarRequest struct {
	AppID          string `json:"appId"` // アプリID
}

func AppStarPost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}
	// fmt.Println(string(data))

	// get user info
	a, _ := oauth.FromContext(ctx)

	// convert request params to struct
	var starReq StarRequest
	if err := json.Unmarshal(data, &starReq); err != nil {
		log.Println("ERROR! json parse", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// get appinfo from db
	appInfo, err := models.AppsInfoTable.FindID(ctx, starReq.AppID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// 重複チェック
	if indexOf(appInfo.StarUsers, a.UserID) == -1 {
		// push the user as starUsers
		appInfo.StarUsers = append(appInfo.StarUsers, a.UserID)
		// update starCount
		appInfo.StarCount = len(appInfo.StarUsers)
	}

	nowTime := time.Now()
	appInfo.UpdateAt = nowTime

	if err := models.AppsInfoTable.Upsert(ctx, appInfo); err != nil {
		log.Println("ERROR! star", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, appInfo.StarUsers)
}

func AppStarDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}
	// fmt.Println(string(data))

	// get user info
	a, _ := oauth.FromContext(ctx)

	// convert request params to struct
	var starReq StarRequest
	if err := json.Unmarshal(data, &starReq); err != nil {
		log.Println("ERROR! json parse", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// get appinfo from db
	appInfo, err := models.AppsInfoTable.FindID(ctx, starReq.AppID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	index := indexOf(appInfo.StarUsers, a.UserID)
	if index != -1 {
		// remove the user from starUsers list
		appInfo.StarUsers = append(appInfo.StarUsers[:index], appInfo.StarUsers[index + 1:]...)
		// update starCount
		appInfo.StarCount = len(appInfo.StarUsers)
	}

	nowTime := time.Now()
	appInfo.UpdateAt = nowTime

	if err := models.AppsInfoTable.Upsert(ctx, appInfo); err != nil {
		log.Println("ERROR! star", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, appInfo.StarUsers)
}
