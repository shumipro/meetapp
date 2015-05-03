package views

import (
    "net/http"

    "github.com/guregu/kami"
    "golang.org/x/net/context"
)

func init() {
    kami.Get("/contact", Constact)
}

func Constact(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    preload := NewHeader(ctx, "Constact", "", "", false)
    ExecuteTemplate(ctx, w, "contact", preload)
}
