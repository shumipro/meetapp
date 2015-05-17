package errors

import (
	"fmt"
	"log"
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
	SendAirbrake(ctx, exception, r)

	//	renderer.JSON(w, 500, "Server Error")
	http.Redirect(w, r, "/error", 302)
}

func SendAirbrake(ctx context.Context, err interface{}, r *http.Request) {
	// send airbrake
	airbrake, ok := goroku.Airbrake(ctx)
	if !ok {
		return
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		airbrake.Notify(err, r)
	}()
}
