package main

import (
	"bytes"
	"io/ioutil"

	"flag"
	"fmt"
	"time"

	"github.com/kyokomi/cloudinary"
	"golang.org/x/net/context"
)

var cURI string

func main() {
	flag.StringVar(&cURI, "uri", "", "")
	flag.Parse()

	ctx := context.Background()
	ctx = cloudinary.NewContext(ctx, cURI)

	// ランダムなファイル名にする
	now := time.Now()
	fileName := fmt.Sprintf("image_%d", now.UnixNano())

	// ファイル読み込み
	data, _ := ioutil.ReadFile("image.png")

	// アップロードする
	cloudinary.UploadStaticImage(ctx, fileName, bytes.NewBuffer(data))

	// アップロードされてるのを確認する
	fmt.Println(cloudinary.ResourceURL(ctx, fileName))
}
