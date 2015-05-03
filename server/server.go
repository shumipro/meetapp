package server

import (
	"log"
	"runtime"

	"net/http"

	"github.com/guregu/kami"
	"github.com/shumipro/meetapp/server/db"
	"github.com/shumipro/meetapp/server/errors"
	"github.com/shumipro/meetapp/server/oauth"
	"github.com/shumipro/meetapp/server/views"
	"golang.org/x/net/context"
)

// Serve start Serve
func Serve() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	ctx := context.Background()
	ctx = db.OpenMongoDB(ctx) // insert mongoDB
	defer db.CloseMongoDB(ctx)
	ctx = db.OpenRedis(ctx) // insert redis
	defer db.CloseRedis(ctx)

	ctx = oauth.WithFacebook(ctx)

	// TODO: とりあえず
	ctx = views.InitTemplates(ctx, "./")

	kami.Context = ctx
	kami.PanicHandler = errors.PanicHandler

	// middleware
	kami.Use("/", oauth.Login)
	kami.Use("/u/", oauth.LoginCheck) // /u以下のpathはloginチェックする

	fileServer := http.FileServer(http.Dir("public"))
	for _, name := range []string{
		"/css/*css",
		"/dist/*dist",
		"/img/*img",
	} {
		kami.Get(name, func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			fileServer.ServeHTTP(w, r)
		})
	}

	log.Println("Starting server...")
	log.Println("GOMAXPROCS: ", cpus)
	kami.Serve()
}
