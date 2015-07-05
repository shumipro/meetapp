package views

import (
	"io/ioutil"
	"log"
	"net/http"

	"html/template"

	"github.com/guregu/kami"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"golang.org/x/net/context"
)

func init() {
	kami.Post("/markdown/preview", MarkdownPreview)
}

func MarkdownPreview(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR!", err)
		renderer.JSON(w, 400, err.Error())
		return
	}

	log.Println(string(data))

	safe := template.HTMLEscapeString(string(data))
	unsafe := blackfriday.MarkdownCommon([]byte(safe))
	markdown := template.HTML(string(bluemonday.UGCPolicy().SanitizeBytes(unsafe)))

	renderer.Data(w, 200, []byte(markdown))
}
