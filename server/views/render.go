package views

import (
	"html/template"
	"path/filepath"

	"fmt"

	"github.com/unrolled/render"
	"golang.org/x/net/context"
)

// TODO: Templateの読み込みあとで別パッケージにする

type TemplateHeader struct {
	Title      string
	NavTitle   string
	SubTitle   string
	ShowBanner bool
}

type templateKey string

var renderer = render.New(render.Options{})

func InitTemplates(ctx context.Context, appRoot string) context.Context {
	path := filepath.FromSlash("views")

	// TODO: あとでディレクトリ指定でいけるようにする

	pageNames := []string{
		"index",
		"login",
		"app/detail",
		"app/list",
		"app/register",
		"error",
		"about",
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
