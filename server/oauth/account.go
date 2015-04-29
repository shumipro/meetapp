package oauth

import (
	"github.com/shumipro/meetapp/server/db"
	"golang.org/x/net/context"
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

func GetAccountByToken(ctx context.Context, authToken string) (Account, error) {
	redisDB := db.Redis(ctx)

	userID, err := redisDB.Get("auth:" + authToken).Result()
	if err != nil {
		return Account{}, err
	}
	return Account{userID, authToken}, nil
}
