package views

import (
	"net/http"

	"github.com/guregu/kami"
	"golang.org/x/net/context"
)

func init() {
	kami.Get("/terms", Terms)
}

func Terms(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := NewHeader(ctx, "Terms", "", "", false)
	ExecuteTemplate(ctx, w, "terms", preload)
}
