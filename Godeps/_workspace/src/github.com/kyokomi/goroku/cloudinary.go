package goroku

import (
	"os"

	"io"

	gocloud "github.com/gotsunami/go-cloudinary"
	"github.com/kyokomi/cloudinary"
	"golang.org/x/net/context"
)

type CloudinaryService struct {
	*gocloud.Service
}

func NewCloudinary(ctx context.Context) context.Context {
	return cloudinary.NewContext(ctx, os.Getenv("CLOUDINARY_URL"))
}

func Cloudinary(ctx context.Context) CloudinaryService {
	c := CloudinaryService{}
	c.Service = cloudinary.FromContext(ctx)
	return c
}

func (c CloudinaryService) UploadStaticImage(path string, fileName string, data io.Reader) error {
	_, err := c.Service.UploadStaticImage(fileName, data, path)
	return err
}

func (c CloudinaryService) Resources() ([]*gocloud.Resource, error) {
	return c.Service.Resources(gocloud.ImageType)
}

func (c CloudinaryService) ResourceURL(fileName string) string {
	return c.Service.Url(fileName, gocloud.ImageType)
}

func (c CloudinaryService) DeleteStaticImage(path string, fileName string) error {
	return c.Service.Delete(fileName, path, gocloud.ImageType)
}
