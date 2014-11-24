package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/GeertJohan/go.rice"
)

var (
	tmplJs *template.Template
	tmplGo *template.Template
)

func setupTemplates() {

	// load rice box
	templatesBox, err := rice.FindBox("templates")
	if err != nil {
		fmt.Printf("Error loading templates: %s\n", err)
		os.Exit(1)
	}

	tmplJs = loadTemplate("ango-service.tmpl.js", templatesBox)
	tmplGo = loadTemplate("ango-service.tmpl.go", templatesBox)
}

func loadTemplate(name string, templatesBox *rice.Box) *template.Template {
	str, err := templatesBox.String(name)
	if err != nil {
		fmt.Printf("Error getting template '%s': %s\n", name, err)
		os.Exit(1)
	}

	tmpl, err := template.New(name).Parse(str)
	if err != nil {
		fmt.Printf("Error parsing template '%s': %s\n", name, err)
		os.Exit(1)
	}
	return tmpl
}
