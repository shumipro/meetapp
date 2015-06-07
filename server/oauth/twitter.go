package oauth

import (
	"fmt"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/mrjones/oauth"
	"golang.org/x/net/context"
)

const twitterAuthCallbackURL = "auth/twitter/callback"

var twitterTokens map[string]*oauth.RequestToken

func init() {
	twitterTokens = make(map[string]*oauth.RequestToken, 0)
}

type TwitterOAuth struct {
	*oauth.Consumer
	callBackURL string
}

func (t TwitterOAuth) GetRequestTokenAndURL() (string, error) {
	token, requestURL, err := t.GetRequestTokenAndUrl(t.callBackURL)
	if token != nil {
		addTwitterToken(token)
	}
	return requestURL, err
}

func WithTwitter(ctx context.Context) context.Context {
	cKey := os.Getenv("TWITTER_CONSUMER_KEY")
	if cKey == "" {
		// TODO:
	}
	cSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	if cSecret == "" {
		// TODO:
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000/"
	}
	callBackURL := baseURL + twitterAuthCallbackURL
	fmt.Println(cKey, cSecret, callBackURL)

	consumer := oauth.NewConsumer(cKey, cSecret, oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	})

	anaconda.SetConsumerKey(cKey)
	anaconda.SetConsumerSecret(cSecret)

	t := TwitterOAuth{consumer, callBackURL}
	return context.WithValue(ctx, authKey("twitter"), t)
}

func Twitter(ctx context.Context) TwitterOAuth {
	conf, _ := ctx.Value(authKey("twitter")).(TwitterOAuth)
	return conf
}

func addTwitterToken(token *oauth.RequestToken) {
	twitterTokens[token.Token] = token
}

func GetTwitterToken(token string) *oauth.RequestToken {
	return twitterTokens[token]
}
