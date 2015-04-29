package oauth

import (
	"net/http"

	"log"

	"golang.org/x/net/context"
)

func Login(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	// Header -> requestParam -> cookieの順番に見に行く
	token := r.Header.Get("Meetup-Auth-Token")
	if token == "" {
		token = r.URL.Query().Get("Meetup-Auth-Token")
		if token == "" {
			ck, err := r.Cookie("Meetup-Auth-Token")
			if err != nil {
				return ctx
			}
			token = ck.Value
		}
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
		http.Redirect(w, r, "/error", 301)
		return nil
	}
	return ctx
}
