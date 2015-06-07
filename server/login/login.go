package login

import (
	"net/http"

	"log"

	"golang.org/x/net/context"
)

const authTokenKey = "Meetapp-Auth-Token"

func Login(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	if a, err := GetAccountBySession(ctx, r); err == nil {
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
