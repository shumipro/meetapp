package views

import (
	"log"
	"net/http"

	"github.com/guregu/kami"
	"golang.org/x/net/context"
)

func init() {
	kami.Get("/", Index)
	kami.Get("/error", Error)
	kami.Get("/about", About)
}

func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := TemplateHeader{
		Title: "MeetApp",
	}
	if err := indexTmpl.Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}

func Error(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := TemplateHeader{
		Title: "Error",
	}
	if err := indexTmpl.Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}

func About(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := TemplateHeader{
		Title: "About",
	}
	if err := indexTmpl.Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err)
		return
	}
}
