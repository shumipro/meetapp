package cloudinary

import (
	"net/url"

	gocloud "github.com/gotsunami/go-cloudinary"
	"golang.org/x/net/context"
)

type key int

const cloudinaryKey key = 0

func NewContext(ctx context.Context, uri string) context.Context {
	cURI, err := url.Parse(uri)
	if err != nil {
		return ctx
	}

	service, err := gocloud.Dial(cURI.String())
	if err != nil {
		return ctx
	}
	return WithCloudinary(ctx, service)
}

func WithCloudinary(ctx context.Context, service *gocloud.Service) context.Context {
	return context.WithValue(ctx, cloudinaryKey, service)
}

func FromContext(ctx context.Context) (*gocloud.Service, bool) {
	c, ok := ctx.Value(cloudinaryKey).(*gocloud.Service)
	return c, ok
}
