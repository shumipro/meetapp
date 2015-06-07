package main

import (
	"errors"
	"flag"
	"log"

	"github.com/kyokomi/cloudinary"
	"golang.org/x/net/context"
)

var cURI string
var dirPath string
var prependPath string

func main() {
	flag.StringVar(&cURI, "uri", "", "")
	flag.StringVar(&dirPath, "dir", "", "")
	flag.StringVar(&prependPath, "prepend", "", "")
	flag.Parse()

	if cURI == "" || dirPath == "" || prependPath == "" {
		log.Fatalln(errors.New("not cloudinary uri"))
	}

	ctx := cloudinary.NewContext(context.Background(), cURI)
	c, ok := cloudinary.FromContext(ctx)
	if !ok {
		log.Fatalln(errors.New("not found cloudinary"))
	}

	_, err := c.UploadStaticImage(dirPath, nil, prependPath)
	if err != nil {
		log.Fatalln("UploadStaticImage:", err)
	}
}
