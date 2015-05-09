package views

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"time"

	"reflect"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/db"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
)

func TestUserProfileUpdate(t *testing.T) {
	ctx := context.Background()
	ctx = db.OpenMongoDB(ctx)
	db.WithMockMongoDB()
	mongoDB := db.MongoDB(ctx)
	defer func() {
		mongoDB.DB(db.MongoDBName()).DropDatabase()
		mongoDB.Close()
	}()

	kami.Use("/", oauth.FakeLogin)
	kami.Use("/u", oauth.LoginCheck)

	kami.Context = ctx

	// Inデータ

	user := models.User{
		ID:            "validUserID",
		ImageURL:      "http://test.png",
		LargeImageURL: "http://large_test.png",
		Name:          "validUserName",
	}

	// Mockデータ投入

	data, err := json.Marshal(user)
	if err != nil {
		t.Errorf("ERROR: %s", err)
	}

	beforeUser := user
	beforeUser.Name = "exampleName"
	beforeUser.ImageURL = "http://example.png"
	beforeUser.LargeImageURL = "http://large_example.png"
	beforeUser.CreateAt = time.Now()
	beforeUser.UpdateAt = time.Now()
	if err := models.UsersTable.Upsert(ctx, beforeUser); err != nil {
		t.Errorf("ERROR: %s", err)
	}

	// テスト実行

	req, err := http.NewRequest("PUT", "/u/api/user", bytes.NewBuffer(data))
	if err != nil {
		t.Errorf("ERROR: %s", err)
	}
	req.Header.Set("Meetapp-Auth-Token", "valid")

	w := httptest.NewRecorder()
	kami.Handler().ServeHTTP(w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())

	// 結果検証

	afterUser, err := readBodyUser(w.Body)
	if err != nil {
		t.Errorf("ERROR: %s", err)
	}

	if !reflect.DeepEqual(user, afterUser) {
		t.Errorf("ERROR: not equal")
	}
}
