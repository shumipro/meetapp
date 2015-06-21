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

	// TODO: update AppInfo in DB with a joined flag?
	// if err := models.AppsInfoTable.Upsert(ctx, appInfo); err != nil {
	// 	panic(err)
	// }

	// Facebook notification API
	// https://developers.facebook.com/docs/graph-api/reference/v2.3/user/notifications
	// for _, m := range appInfo.Members {
	// 	u, _ := models.UsersTable.FindID(ctx, m.UserID)
	// 	if u.FBUser.ID == "" {
	// 		continue
	// 	}
	// 	// check FB account
	// 	urlTmpl := "https://graph.facebook.com/v2.3/{user-id}/notifications?access_token={access-token}&href={href}&template={message}"
	// 	// TODO: get access token and href
	// 	r := strings.NewReplacer("{access-token}", "", "{user-id}", u.FBUser.ID, "{href}", "%2Ftesturl%3Fparam1%3Dvalue1", "{message}", "This+is+a+test+message")
	// 	fbNotificationURL := r.Replace(urlTmpl)
	// 	// send notificatinon
	// 	log.Println(fbNotificationURL)
	// 	req, err := http.NewRequest("POST", fbNotificationURL, nil)
	// 	if err != nil {
	// 		log.Println("ERROR!", err)
	// 		renderer.JSON(w, 400, err.Error())
	// 		return
	// 	}
	// 	client := &http.Client{}
	// 	res, err := client.Do(req)
	// 	if err != nil {
	// 		log.Println("ERROR!", err)
	// 		renderer.JSON(w, 400, err.Error())
	// 		return
	// 	}
	// 	defer res.Body.Close()
	// }

	renderer.JSON(w, 200, appInfo)
}
