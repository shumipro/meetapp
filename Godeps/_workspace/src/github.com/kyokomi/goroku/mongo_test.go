package goroku

import (
	"testing"

	"os"
)

func TestGetHerokuMongoURI(t *testing.T) {
	mongoURI := "mongodb://dbuser:dbpassword@xxxxxxx.mongolab.com:61371/xxxxxxxxxxx13694st"
	os.Setenv("MONGOLAB_URI", mongoURI)

	uri, dbName := getHerokuMongoURI()
	if uri != mongoURI {
		t.Errorf("ERROR: uri %s != %s", uri, mongoURI)
	}

	if dbName != "xxxxxxxxxxx13694st" {
		t.Errorf("ERROR: dbName %s != %s", dbName, "xxxxxxxxxxx13694st")
	}
}
