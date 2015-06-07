package views

import (
	"net/http"

	"github.com/guregu/kami"
	"golang.org/x/net/context"
)

func init() {
	kami.Get("/privacy", Privacy)
}

func Privacy(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	preload := NewHeader(ctx, "Privacy", "", "", false, "", "")
	ExecuteTemplate(ctx, w, r, "privacy", preload)
}
