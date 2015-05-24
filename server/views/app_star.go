package views

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/notification"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
)

func init() {
	// Star API
	kami.Post("/u/api/app/star/:id", APIAppStared)
	kami.Delete("/u/api/app/star/:id", APIAppStarDelete)
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

	user, err := models.UsersTable.FindID(ctx, a.UserID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	notification.SendStar(ctx, user, appInfo)

	if err := models.AppsInfoTable.Upsert(ctx, appInfo); err != nil {
		panic(err)
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
		panic(err)
	}

	renderer.JSON(w, 200, appInfo.StarUsers)
}
