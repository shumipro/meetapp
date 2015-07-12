package errors

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/guregu/kami"
	"github.com/kyokomi/goroku"
	"github.com/unrolled/render"
	"golang.org/x/net/context"
)

var renderer = render.New(render.Options{})

// PanicHandler 500 error
func PanicHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	exception := kami.Exception(ctx)

	fmt.Println("ERROR:", exception)
	fmt.Println(string(debug.Stack()))

	// send airbrake
	if airbrake, ok := goroku.Airbrake(ctx); ok {
		airbrake.Notify(exception, r)
	}

	//	renderer.JSON(w, 500, "Server Error")
	http.Redirect(w, r, "/error", 302)
}
