package oauth

import (
	"net/http"
	"time"
)

const (
	cookieKey = "Meetup-Auth-Token"
)

func readCookieAuthToken(r *http.Request) string {
	ck, err := r.Cookie(cookieKey)
	if err != nil {
		return ""
	}
	return ck.Value
}

func writeCookieAuthToken(w http.ResponseWriter, authToken string, expiry time.Time) {
	var cookie http.Cookie
	cookie.Path = "/"
	cookie.Name = cookieKey
	cookie.Expires = expiry
	cookie.Value = authToken
	http.SetCookie(w, &cookie)
}

func removeCookieAuthToken(w http.ResponseWriter) {
	var cookie http.Cookie
	cookie.Path = "/"
	cookie.Name = cookieKey
	http.SetCookie(w, &cookie)
}
