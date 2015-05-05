package main

import (
	"log"

	"github.com/shumipro/meetapp/server/db"
	"github.com/shumipro/meetapp/server/models"
	"golang.org/x/net/context"
)

var mockUsers = []models.User{
	{ID: "test1", Name: "Satou Yokoyama", FBUser: models.FacebookUser{ID: "facebook1", Name: "Satou Yokoyama"}},
	{ID: "test2", Name: "Yamada Koji", FBUser: models.FacebookUser{ID: "facebook2", Name: "Yamada Koji"}},
}

func main() {
	ctx := context.Background()
	ctx = db.OpenMongoDB(ctx) // insert mongoDB
	defer db.CloseMongoDB(ctx)

	for _, user := range mockUsers {
		if err := models.UsersTable.Upsert(ctx, user); err != nil {
			log.Println(err)
		} else {
			log.Println("OK")
		}
	}
}
