package views

import (
	"net/http"

	"github.com/guregu/kami"
	"golang.org/x/net/context"
)

func init() {
	kami.Get("/about", About)
}

func About(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := NewHeader(ctx, "About", "", "", false)
	ExecuteTemplate(ctx, w, r, "about", preload)
}
