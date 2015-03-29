package server

import (
	"log"
	"runtime"

	"net/http"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/views"
	"golang.org/x/net/context"
)

// Serve start Serve
func Serve() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	ctx := context.Background()
	//	ctx = db.NewContext(ctx)      // insert db
	//	ctx = session.NewContext(ctx) // insert db

	// TODO: とりあえず
	ctx = views.InitTemplates(ctx, "./")

	kami.Context = ctx

	fileServer := http.FileServer(http.Dir("public"))
	for _, name := range []string{
		"/css/*css",
		"/dist/*dist",
		"/html/*html",
		"/img/*img",
		"/js/*js",
		"/stylus/*stylus",
	} {
		kami.Get(name, func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			fileServer.ServeHTTP(w, r)
		})
	}

	log.Println("Starting server...")
	log.Println("GOMAXPROCS: ", cpus)
	kami.Serve()
}
