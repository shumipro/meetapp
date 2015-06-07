package views

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/go-xweb/uuid"
	"github.com/guregu/kami"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"

	"github.com/shumipro/meetapp/server/facebook"
	"github.com/shumipro/meetapp/server/login"
	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
)

func init() {
	kami.Get("/login", Login)
	kami.Get("/logout", Logout)
	kami.Get("/login/facebook", LoginFacebook)
	kami.Get("/login/twitter", LoginTwitter)
	kami.Get("/auth/callback", AuthFacebookCallback) // TODO: Deprecated
	kami.Get("/auth/facebook/callback", AuthFacebookCallback)
	kami.Get("/auth/twitter/callback", AuthTwitterCallback)
}

func Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if _, ok := login.FromContext(ctx); ok {
		// login済みならmypageへ
		http.Redirect(w, r, "/u/mypage", 302)
		return
	}

	preload := NewHeader(ctx, "Login", "", "", false, "", "")
	ExecuteTemplate(ctx, w, r, "login", preload)
}

func Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	login.ResetCacheAuthToken(ctx, w, r)
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

func AuthFacebookCallback(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := oauth.GetFacebookAuthToken(ctx, code)
	if err != nil {
		log.Println("[ERROR] GetFacebookAuthToken", err)
		http.Redirect(w, r, "/error", 302)
		return
	}

	facebookID, res, err := facebook.GetFacebookMe(ctx, token.AccessToken)
	if err != nil {
		panic(err.Error())
	}

	var fbUser models.FacebookUser
	data, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &fbUser); err != nil {
		panic(err)
	}

	var a login.Account
	var user models.User
	var userAuth models.UserAuth
	a, err = login.GetAccountBySession(ctx, r)
	if err == nil {
		// Twitterもしくはfacebookですでにでログイン済み
		user, err = models.UsersTable.FindID(ctx, a.UserID)
	} else {
		user, err = models.UsersTable.FindByFacebookID(ctx, facebookID)
	}

	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	} else if err == mgo.ErrNotFound {
		fmt.Println("新規登録")
		// 新規登録
		user = registerUser(fbUser.Name)
		userAuth.UserID = user.ID
	} else {
		userAuth, _ = models.UserAuthTable.FindID(ctx, user.ID)
		// 登録済み更新
		log.Println("とうろくずみ")

		if userAuth.FacebookToken != "" && user.FBUser.ID != facebookID {
			// 不正ログイン？
			panic(fmt.Errorf("bad login facebookID [%s] != [%s]", user.FBUser.ID, facebookID))
		}
	}

	user.FBUser = fbUser
	userAuth.FacebookToken = token.AccessToken

	if err := models.UsersTable.Upsert(ctx, user); err != nil {
		panic(err)
	}

	if err := models.UserAuthTable.Upsert(ctx, userAuth); err != nil {
		panic(err)
	}

	// RedisでCacheとCookieに書き込む
	if err := login.CacheLoginAccount(ctx, w, r, user.ID); err != nil {
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

	twCli := anaconda.NewTwitterApi(accessToken.Token, accessToken.Secret)
	twUser, err := twCli.GetUsersShow(accessToken.AdditionalData["screen_name"], nil)
	if err != nil {
		panic(err)
	}

	var a login.Account
	var user models.User
	var userAuth models.UserAuth
	a, err = login.GetAccountBySession(ctx, r)
	if err == nil {
		// Twitterもしくはfacebookですでにでログイン済み
		user, err = models.UsersTable.FindID(ctx, a.UserID)
	} else {
		user, err = models.UsersTable.FindByTwitterID(ctx, twUser.Id)
	}
	if err != nil && err != mgo.ErrNotFound {
		panic(err)
	} else if err == mgo.ErrNotFound {
		// 新規登録
		user = registerUser(twUser.Name)
		userAuth.UserID = user.ID
	} else {
		userAuth, _ = models.UserAuthTable.FindID(ctx, user.ID)
		// 登録済み更新
		log.Println("とうろくずみ", user)

		if userAuth.TwitterToken != "" && user.TwitterUser.Id != twUser.Id {
			// 不正ログイン？
			panic(fmt.Errorf("bad login twitter %d != %d", user.TwitterUser.Id, twUser.Id))
		}
	}

	user.TwitterUser = twUser
	userAuth.TwitterToken = accessToken.Token

	if err := models.UsersTable.Upsert(ctx, user); err != nil {
		panic(err)
	}

	if err := models.UserAuthTable.Upsert(ctx, userAuth); err != nil {
		panic(err)
	}

	// RedisでCacheとCookieに書き込む
	if err := login.CacheLoginAccount(ctx, w, r, user.ID); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/u/mypage", 302)
}

func registerUser(name string) models.User {
	user := models.User{}
	user.ID = uuid.New()

	user.Name = name
	user.ImageURL = user.IconImageURL()
	user.LargeImageURL = user.IconLargeImageURL()

	nowTime := time.Now()
	user.CreateAt = nowTime
	user.UpdateAt = nowTime

	return user
}
