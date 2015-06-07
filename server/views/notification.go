package views

import (
	"net/http"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/login"
	"github.com/shumipro/meetapp/server/models"
	"golang.org/x/net/context"
)

func init() {
	// Discussion API
	kami.Get("/u/api/notification", APINotifications)
	kami.Put("/u/api/notification/done", APIAllNotificationsRead)
}

func APINotifications(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := login.FromContext(ctx)

	notification := models.NotificationTable.MustFindID(ctx, a.UserID)
	renderer.JSON(w, 200, notification)
}

func APIAllNotificationsRead(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := login.FromContext(ctx)

	notification := models.NotificationTable.MustFindID(ctx, a.UserID)

	// とりあえず全部Readにする
	if len(notification.Notifications) > 0 {
		for idx, _ := range notification.Notifications {
			notification.Notifications[idx].IsRead = true
		}

		// 10件以上残さない
		notification.TrimNotification(10)

		if err := models.NotificationTable.Upsert(ctx, notification); err != nil {
			panic(err)
		}
	}

	renderer.JSON(w, 200, notification)
}
