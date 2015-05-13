package views

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-xweb/uuid"
	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

func init() {
	kami.Get("/login", Login)
	kami.Get("/logout", Logout)
	kami.Get("/login/facebook", LoginFacebook)
	kami.Get("/auth/callback", AuthCallback)
}

func Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if _, ok := oauth.FromContext(ctx); ok {
		// login済みならmypageへ
		http.Redirect(w, r, "/u/mypage", 302)
		return
	}

	preload := NewHeader(ctx, "Login", "", "", false)
	ExecuteTemplate(ctx, w, r, "login", preload)
}

func Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	oauth.ResetCacheAuthToken(ctx, w, r)
	http.Redirect(w, r, "/login", 302)
}

func LoginFacebook(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c := oauth.Facebook(ctx)
	http.Redirect(w, r, c.AuthCodeURL(""), 302)
}

func AuthCallback(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := oauth.GetFacebookAuthToken(ctx, code)
	if err != nil {
		log.Println("[ERROR] GetFacebookAuthToken", err)
		http.Redirect(w, r, "/error", 302)
		return
	}

	facebookID, res, err := oauth.GetFacebookMe(ctx, token.AccessToken)
	if err != nil {
		panic(err.Error())
	}

	user, err := models.UsersTable.FindByFacebookID(ctx, facebookID)
	if err == mgo.ErrNotFound {
		// 新規
		userID := uuid.New()

		user = models.User{}
		user.ID = userID

		var fbUser models.FacebookUser
		data, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal(data, &fbUser); err != nil {
			panic(err)
		}
		user.Name = fbUser.Name // TODO: 一旦Facebookオンリーなので
		user.FBUser = fbUser
		user.ImageURL = user.IconImageURL()
		user.LargeImageURL = user.IconLargeImageURL()

		nowTime := time.Now()
		user.CreateAt = nowTime
		user.UpdateAt = nowTime

		// 登録する
		if err := models.UsersTable.Upsert(ctx, user); err != nil {
			panic(err)
		} else {
			log.Println("とうろくした")
		}
	} else if err != nil {
		panic(err)
	} else {
		log.Println("とうろくずみ")
	}

	// RedisでCacheとCookieに書き込む
	err = oauth.CacheAuthToken(ctx, w, r, user.ID, *token)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/u/mypage", 302)
}
