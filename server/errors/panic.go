package errors

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/guregu/kami"
	"github.com/kyokomi/goroku"
	"golang.org/x/net/context"
)

// PanicHandler 500 error
func PanicHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	exception := kami.Exception(ctx)

	fmt.Println("ERROR:", exception)
	fmt.Println(string(debug.Stack()))

	// send airbrake
	if n, ok := goroku.Airbrake(ctx); ok {
		notice := n.Notice(exception, r, 3)
		if err := n.SendNotice(notice); err != nil {
			fmt.Println("gobrake failed (%s) reporting error: %v", err, exception)
		}
	}

	//	renderer.JSON(w, 500, "Server Error")
	http.Redirect(w, r, "/error", 302)
}
