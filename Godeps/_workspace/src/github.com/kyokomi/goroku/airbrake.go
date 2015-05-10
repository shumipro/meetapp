package goroku

import (
	"os"

	"golang.org/x/net/context"
	"gopkg.in/airbrake/gobrake.v1"
	"strconv"
)

type airbrakeKey string

func NewAirbrake(ctx context.Context, env string) context.Context {
	projectID, apiKey := os.Getenv("AIRBRAKE_PROJECT_ID"), os.Getenv("AIRBRAKE_API_KEY")
	if projectID == "" || apiKey == "" {
		return ctx
	}

	pid, _ := strconv.ParseInt(projectID, 10, 64)
	airbrake := gobrake.NewNotifier(pid, apiKey)
	airbrake.SetContext("environment", env)
	return context.WithValue(ctx, airbrakeKey("heroku"), airbrake)
}

func Airbrake(ctx context.Context) *gobrake.Notifier {
	return ctx.Value(airbrakeKey("heroku")).(*gobrake.Notifier)
}
