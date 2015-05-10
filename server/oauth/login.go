package oauth

import (
	"net/http"

	"log"

	"golang.org/x/net/context"
)

const authTokenKey = "Meetapp-Auth-Token"

func Login(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	// Header -> requestParam -> cookieの順番に見に行く
	token := r.Header.Get(authTokenKey)
	if token == "" {
		token = r.URL.Query().Get(authTokenKey)
		if token == "" {
			token = readCookieAuthToken(r)
		}
	}
	if token == "" {
		return ctx
	}

	if a, err := GetAccountByToken(ctx, token); err == nil {
		ctx = NewContext(ctx, a)
	}
	return ctx
}

func LoginCheck(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	_, ok := FromContext(ctx)
	if !ok {
		log.Println("[ERROR] Login Error 401")
		http.Redirect(w, r, "/login", 302)
		return nil
	}
	return ctx
}

func FakeLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	token := r.Header.Get(authTokenKey)
	if token == "valid" {
		return NewContext(ctx, Account{UserID: "validUserID", AuthToken: token})
	}
	return ctx
}
