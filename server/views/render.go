package views

import (
	"html/template"
	"path/filepath"

	"fmt"

	"net/http"

	"strings"

	"github.com/shumipro/meetapp/server/models"
	"github.com/shumipro/meetapp/server/oauth"
	"github.com/unrolled/render"
	"golang.org/x/net/context"
)

type Config struct {
	User models.User `json:"user"`
}

type TemplateHeader struct {
	Title      string
	NavTitle   string
	SubTitle   string
	ShowBanner bool
	Config     Config           `json:"config"`
	Constants  models.Constants `json:"constants"`
}

func (t TemplateHeader) EscapeNewline(text string) template.HTML {
	safe := template.HTMLEscapeString(text)
	safe = strings.Replace(safe, "\n", "<br />", -1)
	return template.HTML(safe)
}

func NewHeader(ctx context.Context, title, navTitle, subTitle string, showBanner bool) TemplateHeader {
	a, _ := oauth.FromContext(ctx)

	// TODO: 毎アクセスでmongoとるの微妙・・・ Serverでcacheしてもよさそう
	user, _ := models.UsersTable.FindID(ctx, a.UserID)

	h := TemplateHeader{}
	h.Config = Config{User: user}
	h.Constants = models.AllConstants()

	h.Title = title
	h.SubTitle = subTitle
	h.NavTitle = navTitle
	h.ShowBanner = showBanner

	return h
}

type templateKey string

var renderer = render.New(render.Options{})

func InitTemplates(ctx context.Context, appRoot string) context.Context {
	path := filepath.FromSlash("views")

	// TODO: あとでディレクトリ指定でいけるようにする

	pageNames := []string{
		"index",
		"login",
		"mypage",
		"app/detail",
		"app/list",
		"app/register",
		"error",
		"about",
		"terms",
		"privacy",
		"contact",
	}
	tmplMap := make(map[string]*template.Template, 0)
	for _, name := range pageNames {
		tmplMap[name] = template.Must(template.ParseFiles(filepath.Join(appRoot, path, name+".html")))
	}

	subNames := []string{
		"app/list",
		"app/register",
		"app/components/tile",
		"app/components/listItem",
		"partials/footer",
		"partials/header",
		"partials/nav",
		"partials/scripts",
	}
	for _, name := range subNames {
		subTemplate := template.Must(template.ParseFiles(filepath.Join(appRoot, path, name+".html")))
		fmt.Printf("Template: %+v\n", subTemplate.Name())
		for _, tmpl := range tmplMap {
			tmpl.AddParseTree(name, subTemplate.Tree)
		}
	}

	return context.WithValue(ctx, templateKey("default"), tmplMap)
}

func FromContextTemplate(ctx context.Context, name string) *template.Template {
	tmpls, ok := ctx.Value(templateKey("default")).(map[string]*template.Template)
	if !ok {
		panic("not template")
	}
	return tmpls[name]
}

func ExecuteTemplate(ctx context.Context, w http.ResponseWriter, name string, data interface{}) {
	if err := FromContextTemplate(ctx, name).Execute(w, data); err != nil {
		executeError(ctx, w, err)
	}
}
