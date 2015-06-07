package oauth

import (
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type keyType int

var accountKey keyType = 0

type Account struct {
	UserID    string
	AuthToken string
}

type responseUserID struct {
	UserID int64 `json:"id"`
}

func NewContext(ctx context.Context, a Account) context.Context {
	return context.WithValue(ctx, accountKey, a)
}

func FromContext(ctx context.Context) (Account, bool) {
	u, ok := ctx.Value(accountKey).(Account)
	return u, ok
}

func GetAccountByToken(ctx context.Context, r *http.Request) (Account, error) {
	return readSessionAuthToken(ctx, r)
}

func CacheAuthToken(ctx context.Context, w http.ResponseWriter, r *http.Request, userID string, token oauth2.Token) error {
	return writeSessionAuthToken(ctx, w, r, Account{userID, token.AccessToken})
}

func ResetCacheAuthToken(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	removeSessionAuthToken(ctx, w, r)
}
