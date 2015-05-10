package views

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/guregu/kami"
	"github.com/kyokomi/cloudinary"
	"github.com/shumipro/meetapp/server/oauth"
	"golang.org/x/net/context"
)

func init() {
	// API
	kami.Post("/u/api/upload/image", UploadImage)
}

type UploadImageResponse struct {
	LargeImageURL string
	ImageURL      string
}

// UploadImage Cloudinaryに画像をUploadして画像のURLを返すAPI
func UploadImage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	a, _ := oauth.FromContext(ctx)

	formFile, _, err := r.FormFile("file")
	if err != nil {
		renderer.JSON(w, 400, err.Error())
		return
	}
	defer formFile.Close()

	// Uploadする
	fileName := fmt.Sprintf("%s_%d", a.UserID, time.Now().UnixNano())
	if err := cloudinary.UploadStaticImage(ctx, fileName, formFile); err != nil {
		panic(err)
	}

	// Uploadした画像のURLを取得する
	res := UploadImageResponse{}
	res.LargeImageURL = cloudinary.ResourceURL(ctx, fileName)
	res.ImageURL = strings.Replace(res.LargeImageURL, "image/upload", "image/upload/w_96,h_96", 1)

	renderer.JSON(w, 200, res)
}
