package views

import (
	"log"
	"net/http"

	"time"

	"github.com/go-xweb/uuid"
	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/db"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func init() {
	kami.Get("/login", Login)
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
	if err := FromContextTemplate(ctx, "login").Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
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

	expiry := token.Expiry.Sub(time.Now())

	userID := uuid.New()
	_, err = redisDB.SetEx("auth:"+token.AccessToken, expiry, userID).Result()
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
