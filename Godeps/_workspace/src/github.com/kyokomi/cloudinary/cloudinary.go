package cloudinary

import (
	"io"

	gocloud "github.com/gotsunami/go-cloudinary"
	"golang.org/x/net/context"
)

func UploadStaticImage(ctx context.Context, fileName string, data io.Reader) error {
	_, err := FromContext(ctx).UploadStaticImage(fileName, data, "")
	return err
}

func Resources(ctx context.Context) ([]*gocloud.Resource, error) {
	return FromContext(ctx).Resources(gocloud.ImageType)
}

func ResourceURL(ctx context.Context, fileName string) string {
	return FromContext(ctx).Url(fileName, gocloud.ImageType)
}
