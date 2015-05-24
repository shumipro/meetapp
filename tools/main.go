package main

import (
	"log"
	"time"

	"github.com/kyokomi/goroku"
	"github.com/shumipro/meetapp/server/models"
	"golang.org/x/net/context"
)

var mockUsers = []models.User{
	{ID: "test1", Name: "Satou Yokoyama", FBUser: models.FacebookUser{ID: "facebook1", Name: "Satou Yokoyama"}},
	{ID: "test2", Name: "Yamada Koji", FBUser: models.FacebookUser{ID: "facebook2", Name: "Yamada Koji"}},
}

var mockNotifications = []models.UserNotification{
	{UserID: "test1", Notifications: []models.Notification{
		{"1", models.NotificationDiscussion, "1", "Messasge", "Hoge", false, time.Now()},
	}},
}

func main() {
	ctx := context.Background()
	ctx = goroku.OpenMongoDB(ctx) // insert mongoDB
	defer goroku.CloseMongoDB(ctx)

	for _, user := range mockUsers {
		if err := models.UsersTable.Upsert(ctx, user); err != nil {
			log.Println(err)
		} else {
			log.Println("OK")
		}
	}

	for _, notification := range mockNotifications {
		if err := models.NotificationTable.Upsert(ctx, notification); err != nil {
			log.Println(err)
		} else {
			log.Println("OK")
		}
	}
}
