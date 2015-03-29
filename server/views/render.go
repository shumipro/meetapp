package views

import (
	"html/template"
	"path/filepath"

	"fmt"

	"github.com/unrolled/render"
)

// TODO: Templateの読み込みあとで別パッケージにする

type TemplateHeader struct {
	Title      string
	SubTitle   string
	ShowBanner bool
}

var renderer = render.New(render.Options{})

var indexTmpl *template.Template
var appListTmpl *template.Template
var appRegisterTmpl *template.Template
var errorTmpl *template.Template
var aboutTmpl *template.Template

func InitTemplates(appRoot string) {
	path := filepath.FromSlash("views")

	subNames := []string{
		"app/list",
		"app/register",
		"partials/footer",
		"partials/header",
		"partials/nav",
		"partials/scripts",
	}

	indexTmpl = template.Must(template.ParseFiles(filepath.Join(appRoot, path, "index.html")))
	appListTmpl = template.Must(template.ParseFiles(filepath.Join(appRoot, path, "app/list.html")))
	appRegisterTmpl = template.Must(template.ParseFiles(filepath.Join(appRoot, path, "app/register.html")))
	errorTmpl = template.Must(template.ParseFiles(filepath.Join(appRoot, path, "error.html")))
	aboutTmpl = template.Must(template.ParseFiles(filepath.Join(appRoot, path, "about.html")))

	for _, name := range subNames {
		subTemplate := template.Must(template.ParseFiles(filepath.Join(appRoot, path, name+".html")))
		fmt.Printf("Template: %+v\n", subTemplate.Name())
		indexTmpl.AddParseTree(name, subTemplate.Tree)
		appListTmpl.AddParseTree(name, subTemplate.Tree)
		appRegisterTmpl.AddParseTree(name, subTemplate.Tree)
		errorTmpl.AddParseTree(name, subTemplate.Tree)
		aboutTmpl.AddParseTree(name, subTemplate.Tree)
	}
}
