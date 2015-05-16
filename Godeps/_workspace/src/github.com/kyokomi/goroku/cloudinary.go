package goroku

import (
	"os"

	"io"

	gocloud "github.com/gotsunami/go-cloudinary"
	"github.com/kyokomi/cloudinary"
	"golang.org/x/net/context"
	"errors"
)

var (
	NotFoundCloudinary = errors.New("not found cloudinary")
)

type CloudinaryService struct {
	*gocloud.Service
}

func Cloudinary(ctx context.Context) CloudinaryService {
	c := CloudinaryService{}
	cl, ok := cloudinary.FromContext(ctx)
	if ok {
		c.Service = cl
	}
	return c
}

func NewCloudinary(ctx context.Context) context.Context {
	return cloudinary.NewContext(ctx, os.Getenv("CLOUDINARY_URL"))
}

func (c CloudinaryService) UploadStaticImage(path string, fileName string, data io.Reader) error {
	if c.Service != nil {
		return NotFoundCloudinary
	}
	_, err := c.Service.UploadStaticImage(fileName, data, path)
	return err
}

func (c CloudinaryService) Resources() ([]*gocloud.Resource, error) {
	if c.Service != nil {
		return nil, NotFoundCloudinary
	}
	return c.Service.Resources(gocloud.ImageType)
}

func (c CloudinaryService) ResourceURL(fileName string) string {
	if c.Service != nil {
		return ""
	}
	return c.Service.Url(fileName, gocloud.ImageType)
}

func (c CloudinaryService) DeleteStaticImage(path string, fileName string) error {
	if c.Service != nil {
		return NotFoundCloudinary
	}
	return c.Service.Delete(fileName, path, gocloud.ImageType)
}
