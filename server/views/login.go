package views

import (
	"net/http"

	"time"

	"encoding/json"
	"fmt"

	"github.com/go-xweb/uuid"
	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"github.com/k0kubun/pp"
	"github.com/ChimeraCoder/anaconda"
)

func init() {
	kami.Get("/login", Login)
	kami.Get("/logout", Logout)
	kami.Get("/login/facebook", LoginFacebook)
	kami.Get("/login/twitter", LoginTwitter)
	kami.Get("/auth/callback", AuthCallback) // TODO: Deprecated
	kami.Get("/auth/facebook/callback", AuthCallback)
	kami.Get("/auth/twitter/callback", AuthTwitterCallback)
}

func Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if _, ok := oauth.FromContext(ctx); ok {
		// login済みならmypageへ
		http.Redirect(w, r, "/u/mypage", 302)
		return
	}

	preload := NewHeader(ctx, "Login", "", "", false)
	ExecuteTemplate(ctx, w, "login", preload)
}

func Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	oauth.ResetCacheAuthToken(ctx, w)
	http.Redirect(w, r, "/login", 302)
}

func LoginFacebook(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c := oauth.Facebook(ctx)
	http.Redirect(w, r, c.AuthCodeURL(""), 302)
}

func LoginTwitter(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c := oauth.Twitter(ctx)

	requestURL, err := c.GetRequestTokenAndURL()
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, requestURL, 302)
}

func AuthCallback(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := oauth.GetFacebookAuthToken(ctx, code)
	if err != nil {
		panic(err.Error())
	}

	facebookID, res, err := oauth.GetFacebookMe(ctx, token.AccessToken)
	if err != nil {
		panic(err.Error())
	}

	// TODO: Twitterでログイン済みの考慮が必要

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
			fmt.Println("とうろくした")
		}
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("とうろくずみ")
	}

	// RedisでCacheとCookieに書き込む
	err = oauth.CacheAuthToken(ctx, w, user.ID, *token)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/u/mypage", 302)
}

func AuthTwitterCallback(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c := oauth.Twitter(ctx)

	values := r.URL.Query()
	verificationCode := values.Get("oauth_verifier")
	tokenKey := values.Get("oauth_token")

	accessToken, err := c.AuthorizeToken(oauth.GetTwitterToken(tokenKey), verificationCode)
	if err != nil {
		panic(err)
	}

	cli := anaconda.NewTwitterApi(accessToken.Token, accessToken.Secret)
	user, err := cli.GetUsersShow(accessToken.AdditionalData["screen_name"], nil)
	if err != nil {
		panic(err)
	}

	// TODO: Facebookでログイン済みを考慮しないといけない

	pp.Println(user)
}
