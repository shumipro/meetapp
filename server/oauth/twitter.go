package oauth
import (
	"github.com/ChimeraCoder/anaconda"
	"os"
	"errors"
	"fmt"
)

type TwClient struct {
	api anaconda.TwitterApi
}

// NewTwitterClient creates a new Twitter client with the given access information.
func NewTwitterClient() (*TwClient, error) {
	ckey, err := getenv("TWITTER_CONSUMER_KEY")
	csec, err := getenv("TWITTER_CONSUMER_SECRET")
	atok, err := getenv("TWITTER_ACCESS_TOKEN")
	asec, err := getenv("TWITTER_ACCESS_SECRET")
	if err != nil {
		return nil, err
	}
	anaconda.SetConsumerKey(ckey)
	anaconda.SetConsumerSecret(csec)
	return &TwClient{*anaconda.NewTwitterApi(atok, asec)}, nil
}

func getenv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return val, errors.New(fmt.Sprintf("The specified OS environment does not exist: %s.", key))
	}
	return val, nil
}

func (c TwClient) Tweet(msg string) (int64, error) {
	tweet, err := c.api.PostTweet(msg, nil)
	if err != nil {
		return 0, err
	} else {
		return tweet.Id, nil
	}
}