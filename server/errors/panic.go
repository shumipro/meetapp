package errors

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/guregu/kami"
	"github.com/unrolled/render"
	"golang.org/x/net/context"
	"github.com/kyokomi/goroku"
)

var renderer = render.New(render.Options{})

// PanicHandler 500 error
func PanicHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	exception := kami.Exception(ctx)

	fmt.Println("ERROR:", exception)
	fmt.Println(string(debug.Stack()))

	// send airbrake
	go goroku.Airbrake(ctx).Notify(exception, r)

	//	renderer.JSON(w, 500, "Server Error")
	http.Redirect(w, r, "/error", 302)
}
