package views

import (
	"log"
	"net/http"

	"github.com/guregu/kami"
	"github.com/kyokomi/goroku"
	"golang.org/x/net/context"
)

func init() {
	kami.Get("/error", Error)
}

func Error(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	executeError(ctx, w, r, nil)
}

func executeError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	preload := NewHeader(ctx, "Error", "", "", false, "", "")

	if err != nil {
		log.Println("ERROR!", err)
		if n, ok := goroku.Airbrake(ctx); ok {
			notice := n.Notice(err, r, 3)
			if sendErr := n.SendNotice(notice); sendErr != nil {
				log.Printf("gobrake failed (%s) reporting error: %v\n", sendErr, err)
			}
		}
		preload.SubTitle = err.Error()
	}

	if err := FromContextTemplate(ctx, "error").Execute(w, preload); err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}
}
