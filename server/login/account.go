package login

import (
	"net/http"

	"golang.org/x/net/context"
)

type keyType int

var accountKey keyType = 0

type Account struct {
	UserID string
}

func NewContext(ctx context.Context, a Account) context.Context {
	return context.WithValue(ctx, accountKey, a)
}

func FromContext(ctx context.Context) (Account, bool) {
	u, ok := ctx.Value(accountKey).(Account)
	return u, ok
}

func GetAccountBySession(ctx context.Context, r *http.Request) (Account, error) {
	return readSessionAuthToken(ctx, r)
}

func CacheLoginAccount(ctx context.Context, w http.ResponseWriter, r *http.Request, userID string) error {
	return writeSessionAuthToken(ctx, w, r, Account{userID})
}

func ResetCacheAuthToken(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	removeSessionAuthToken(ctx, w, r)
}
