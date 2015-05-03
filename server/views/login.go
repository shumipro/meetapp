package views

import (
	"net/http"

	"time"

	"encoding/json"
	"fmt"

	"github.com/go-xweb/uuid"
	"github.com/guregu/kami"
	fb "github.com/huandu/facebook"
	"github.com/shumipro/meetapp/server/db"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
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
		http.Redirect(w, r, "/u/mypage", 301)
		return
	}

	preload := TemplateHeader{
		Title: "Login",
	}
	ExecuteTemplate(ctx, w, "login", preload)
}

func Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := oauth.FromContext(ctx)
	redisDB := db.Redis(ctx)
	redisDB.Del("auth:" + a.AuthToken)
	removeCookieAuthToken(w)

	http.Redirect(w, r, "/login", 301)
}

func LoginFacebook(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c := oauth.Facebook(ctx)
	http.Redirect(w, r, c.AuthCodeURL(""), 301)
}

func AuthCallback(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c := oauth.Facebook(ctx)
	redisDB := db.Redis(ctx)

	code := r.FormValue("code")
	token, err := c.Exchange(oauth2.NoContext, code)
	if err != nil {
		panic(err)
	}

	// TODO: kyokomi あとでリファクタします...

	// Redisから登録済みかを取得
	var user models.User
	res, err := fb.Get("/me", fb.Params{
		"access_token": token.AccessToken,
	})
	if err != nil {
		panic(err)
	}

	facebookID := res["id"].(string)
	user, err = models.UsersTable().FindByFacebookID(ctx, facebookID)
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
		user.ImageURL = user.IconImageURL()
		user.FBUser = fbUser

		nowTime := time.Now()
		user.CreateAt = nowTime
		user.UpdateAt = nowTime

		// 登録する
		if err := models.UsersTable().Upsert(ctx, user); err != nil {
			panic(err)
		} else {
			fmt.Println("とうろくした")
		}
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("とうろくずみ")
	}

	// RedisでCacheしてる
	expiry := token.Expiry.Sub(time.Now())
	_, err = redisDB.SetEx("auth:"+token.AccessToken, expiry, user.ID).Result()
	if err != nil {
		panic(err)
	}
	writeCookieAuthToken(w, token.AccessToken, token.Expiry)

	http.Redirect(w, r, "/u/mypage", 301)
}

func writeCookieAuthToken(w http.ResponseWriter, authToken string, expiry time.Time) {
	// TODO: とりあえずCookieに焼く
	var cookie http.Cookie
	cookie.Path = "/"
	cookie.Name = "Meetup-Auth-Token"
	cookie.Expires = expiry
	cookie.Value = authToken
	http.SetCookie(w, &cookie)
}

func removeCookieAuthToken(w http.ResponseWriter) {
	// TODO: とりあえずCookieに焼く
	var cookie http.Cookie
	cookie.Path = "/"
	cookie.Name = "Meetup-Auth-Token"
	http.SetCookie(w, &cookie)
}
