package oauth

import (
	"log"
	"net/http"
	"time"

	"github.com/shumipro/meetapp/server/db"
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

func GetAccountByToken(ctx context.Context, authToken string) (Account, error) {
	redisDB := db.Redis(ctx)

	userID, err := redisDB.Get("auth:" + authToken).Result()
	if err != nil {
		return Account{}, err
	}
	return Account{userID, authToken}, nil
}

func CacheAuthToken(ctx context.Context, w http.ResponseWriter, userID string, token oauth2.Token) error {
	redisDB := db.Redis(ctx)

	_, err := redisDB.SetEx("auth:"+token.AccessToken, token.Expiry.Sub(time.Now()), userID).Result()
	if err != nil {
		log.Println("ERROR: Redis.SetEx", err, token, userID)
		return err
	}
	writeCookieAuthToken(w, token.AccessToken, token.Expiry)

	return nil
}

func ResetCacheAuthToken(ctx context.Context, w http.ResponseWriter) {
	a, _ := FromContext(ctx)
	redisDB := db.Redis(ctx)

	redisDB.Del("auth:" + a.AuthToken)
	removeCookieAuthToken(w)
}
