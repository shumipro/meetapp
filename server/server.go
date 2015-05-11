package server

import (
	"log"
	"runtime"

	"net/http"

	"github.com/guregu/kami"
	"github.com/kyokomi/goroku"
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
	ctx = goroku.OpenMongoDB(ctx) // insert mongoDB
	defer goroku.CloseMongoDB(ctx)
	ctx = goroku.OpenRedis(ctx) // insert redis
	defer goroku.CloseRedis(ctx)

	ctx = oauth.WithFacebook(ctx)
	ctx = goroku.NewCloudinary(ctx)

	ctx = goroku.NewAirbrake(ctx, "production")

	ctx = oauth.NewSessionStore(ctx)
	defer oauth.CloseSessionStore(ctx)

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
		"/favicon.ico",
		"/robots.txt",
		"/sitemap.xml",
	} {
		kami.Get(name, func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			fileServer.ServeHTTP(w, r)
		})
	}

	log.Println("Starting server...")
	log.Println("GOMAXPROCS: ", cpus)
	kami.Serve()
}
