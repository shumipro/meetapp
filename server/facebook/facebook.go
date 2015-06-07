package facebook

import (
	"github.com/huandu/facebook"
	"golang.org/x/net/context"
)

func GetFacebookMe(ctx context.Context, authToken string) (facebookID string, res facebook.Result, err error) {
	res, err = facebook.Get("/me", facebook.Params{
		"access_token": authToken,
	})
	if err != nil {
		return "", nil, err
	}

	facebookID = res["id"].(string)

	return facebookID, res, nil
}
