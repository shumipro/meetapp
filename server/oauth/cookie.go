package oauth

import (
	"log"
	"os"

	"errors"
	"net/http"

	"github.com/boj/redistore"
	"github.com/kyokomi/goroku"
	"golang.org/x/net/context"
)

type sessionStoreKey string

const (
	cookieKey       = "meetapp-secret"
	idleConnections = 10
)

var logger = log.New(os.Stderr, "ERROR", log.Llongfile|log.LstdFlags)

func NewSessionStore(ctx context.Context) context.Context {
	addr, password := goroku.GetHerokuRedisAddr()
	store, err := redistore.NewRediStore(idleConnections, "tcp", addr, password, []byte("test"))
	if err != nil {
		logger.Println(err)
		return ctx
	}
	return withSessionStore(ctx, store)
}

func withSessionStore(ctx context.Context, store *redistore.RediStore) context.Context {
	return context.WithValue(ctx, sessionStoreKey(cookieKey), store)
}

func FromSessionStore(ctx context.Context) *redistore.RediStore {
	return ctx.Value(sessionStoreKey(cookieKey)).(*redistore.RediStore)
}

func CloseSessionStore(ctx context.Context) error {
	return FromSessionStore(ctx).Close()
}

func readSessionAuthToken(ctx context.Context, r *http.Request) (Account, error) {
	session, err := FromSessionStore(ctx).Get(r, cookieKey)
	if err != nil {
		return Account{}, err
	}

	userID, _ := session.Values["UserID"].(string)
	authToken, _ := session.Values["AuthToken"].(string)
	if userID == "" || authToken == "" {
		return Account{}, errors.New("not login")
	}

	return Account{userID, authToken}, nil
}

func writeSessionAuthToken(ctx context.Context, w http.ResponseWriter, r *http.Request, account Account) error {
	session, err := FromSessionStore(ctx).Get(r, cookieKey)
	if err != nil {
		log.Println(err)
		return err
	}
	session.Values["UserID"] = account.UserID
	session.Values["AuthToken"] = account.AuthToken
	if err := session.Save(r, w); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func removeSessionAuthToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	session, err := FromSessionStore(ctx).Get(r, cookieKey)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
