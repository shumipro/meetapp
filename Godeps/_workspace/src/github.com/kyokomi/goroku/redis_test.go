package goroku

import (
	"os"
	"testing"
)

func TestGetHerokuRedisAddr(t *testing.T) {
	os.Setenv("REDISTOGO_URL", "redis://redistogo:<password>@mummichog.redistogo.com:11068/")

	addr, password := getHerokuRedisAddr()
	if addr != "mummichog.redistogo.com:11068" {
		t.Errorf("ERROR: addr %s != %s", addr, "mummichog.redistogo.com:11068")
	}

	if password != "<password>" {
		t.Errorf("ERROR: password %s != %s", password, "<password>")
	}
}
