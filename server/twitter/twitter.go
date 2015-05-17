package twitter

import (
	"errors"
	"fmt"
	"os"

	"log"

	"github.com/ChimeraCoder/anaconda"
	"golang.org/x/net/context"
)

type key string

// NewContext creates a new Twitter client with the given access information.
func NewContext(ctx context.Context) context.Context {
	t, err := newTwitterClient()
	if err != nil {
		log.Println(err)
		return ctx
	}
	return context.WithValue(ctx, key("twClient"), t)
}

func FromContext(ctx context.Context) (TwClient, bool) {
	t, ok := ctx.Value(key("twClient")).(TwClient)
	return t, ok
}

type TwClient struct {
	*anaconda.TwitterApi
}

func newTwitterClient() (TwClient, error) {
	var err error
	getEnvFunc := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			err = errors.New(fmt.Sprintf("The specified OS environment does not exist: %s.", key))
		}
		return val
	}

	ckey := getEnvFunc("TWITTER_CONSUMER_KEY")
	csec := getEnvFunc("TWITTER_CONSUMER_SECRET")
	atok := getEnvFunc("TWITTER_ACCESS_TOKEN")
	asec := getEnvFunc("TWITTER_ACCESS_SECRET")
	if err != nil {
		return TwClient{}, err
	}
	anaconda.SetConsumerKey(ckey)
	anaconda.SetConsumerSecret(csec)

	t := TwClient{}
	t.TwitterApi = anaconda.NewTwitterApi(atok, asec)
	return t, nil
}

func (c TwClient) Tweet(msg string) (int64, error) {
	tweet, err := c.PostTweet(msg, nil)
	if err != nil {
		return 0, err
	}

	return tweet.Id, nil
}
