package views

import (
	"log"
	"net/http"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/login"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/notification"
	"golang.org/x/net/context"
)

func init() {
	kami.Post("/u/api/app/join/:id", APIAppJoin)
}

func APIAppJoin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := login.FromContext(ctx)
	appID := kami.Param(ctx, "id")

	appInfo, err := models.AppsInfoTable.FindID(ctx, appID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	user, err := models.UsersTable.FindID(ctx, a.UserID)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	notification.SendJoin(ctx, user, appInfo)

	// TODO: set AppInfo to DB?
	// if err := models.AppsInfoTable.Upsert(ctx, appInfo); err != nil {
	// 	panic(err)
	// }

	renderer.JSON(w, 200, appInfo)
}
