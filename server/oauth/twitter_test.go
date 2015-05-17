package oauth
import (
	"testing"
)

// This test actually posts a tweet so make sure it does not run with "go test" by default.
func _TestTweet(t *testing.T) {
	client, err := NewTwitterClient()
	if err != nil {
		t.Errorf("Failed to initialize Twitter client: %s", err)
	}
	id, err := client.Tweet("This is a test posted by a program.")
	if err != nil {
		t.Errorf("Failed to tweet: %s", err)
		t.Fail()
	}
	if id == 0 {
		t.Errorf("Twitter ID is 0")
		t.Fail()
	}
}