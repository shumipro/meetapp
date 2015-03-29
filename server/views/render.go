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
	SubTitle   string
	ShowBanner bool
}

type templateKey string

var renderer = render.New(render.Options{})

func InitTemplates(ctx context.Context, appRoot string) context.Context {
	path := filepath.FromSlash("views")

	subNames := []string{
		"app/list",
		"app/register",
		"partials/footer",
		"partials/header",
		"partials/nav",
		"partials/scripts",
	}

	tmplMap := make(map[string]*template.Template, 0)
	{
		tmplMap["index"] = template.Must(template.ParseFiles(filepath.Join(appRoot, path, "index.html")))

		tmplMap["app/list"] = template.Must(template.ParseFiles(filepath.Join(appRoot, path, "app/list.html")))
		tmplMap["app/register"] = template.Must(template.ParseFiles(filepath.Join(appRoot, path, "app/register.html")))

		tmplMap["error"] = template.Must(template.ParseFiles(filepath.Join(appRoot, path, "error.html")))
		tmplMap["about"] = template.Must(template.ParseFiles(filepath.Join(appRoot, path, "about.html")))
	}

	for _, name := range subNames {
		subTemplate := template.Must(template.ParseFiles(filepath.Join(appRoot, path, name+".html")))
		fmt.Printf("Template: %+v\n", subTemplate.Name())
		for _, tmpl := range tmplMap {
			tmpl.AddParseTree(name, subTemplate.Tree)
		}
	}

	return context.WithValue(ctx, templateKey(""), tmplMap)
}

func FromContextTemplate(ctx context.Context, name string) *template.Template {
	tmpls, ok := ctx.Value(templateKey("")).(map[string]*template.Template)
	if !ok {
		panic("not template")
	}
	return tmpls[name]
}
