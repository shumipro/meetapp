package views

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/notification"
	"golang.org/x/net/context"
	"html/template"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"
)

func init() {
	// Discussion API
	kami.Post("/u/api/app/discussion", APIAppDiscussion)
	kami.Delete("/u/api/app/discussion/:id", APIAppDiscussionDelete)
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
	discussionReq.DiscussionInfo.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
	discussionReq.DiscussionInfo.Timestamp = nowTime

	// Message To Markdown
	message := discussionReq.DiscussionInfo.Message
	if message != "" {
		safe := template.HTMLEscapeString(message)
		unsafe := blackfriday.MarkdownCommon([]byte(safe))
		discussionReq.DiscussionInfo.MessageMD = template.HTML(string(bluemonday.UGCPolicy().SanitizeBytes(unsafe)))
	}

	// push a discussionInfo
	appInfo.Discussions = append(appInfo.Discussions, discussionReq.DiscussionInfo)
	appInfo.UpdateAt = nowTime

	if err := models.AppsInfoTable.Upsert(ctx, appInfo); err != nil {
		panic(err)
	}

	notification.SendDiscussion(ctx, discussionReq.DiscussionInfo, appInfo)

	renderer.JSON(w, 200, appInfo.Discussions)
}

func APIAppDiscussionDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	discussionID := kami.Param(ctx, "id")

	if err := models.AppsInfoTable.DeleteDiscussionByID(ctx, discussionID); err != nil {
		log.Println("ERROR! discussion", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, discussionID)
}
