package views

import (
	"log"
	"net/http"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

func init() {
	// Discussion API
	kami.Get("/u/api/notification", APINotifications)
	kami.Put("/u/api/notification/done", APIAllNotificationsRead)
}

func APINotifications(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := oauth.FromContext(ctx)

	notification, err := models.NotificationTable.FindID(ctx, a.UserID)
	if err == mgo.ErrNotFound {
		// 空の時は代わりを作ってあげる
		notification = models.UserNotification{}
		notification.UserID = a.UserID
		notification.Notifications = []models.Notification{}
	} else if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	renderer.JSON(w, 200, notification)
}

func APIAllNotificationsRead(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := oauth.FromContext(ctx)

	notification, err := models.NotificationTable.FindID(ctx, a.UserID)
	if err != nil {
		log.Println("ERROR!", err, a.UserID)
		renderer.JSON(w, 400, err.Error())
		return
	}

	// とりあえず全部Readにする
	for idx, _ := range notification.Notifications {
		notification.Notifications[idx].IsRead = true
	}

	// 10件以上残さない
	notification.TrimNotification(10)

	if err := models.NotificationTable.Upsert(ctx, notification); err != nil {
		panic(err)
	}

	renderer.JSON(w, 200, notification)
}
