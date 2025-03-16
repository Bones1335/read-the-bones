package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/gomarkdown/markdown"
)

func mdToHTML() string {
	mdFile, err := os.ReadFile("content/connemara_Mathieu_Nicolas.md")
	if err != nil {
		fmt.Printf("error reading markdown file: %v", err)
		return ""
	}

	html := markdown.ToHTML(mdFile, nil, nil)

	return string(html)
}

func handlerBlog(w http.ResponseWriter, r *http.Request) {
	md := mdToHTML()

	temp, err := template.ParseFiles("templates/layout.html", "templates/blog.html")
	if err != nil {
		fmt.Printf("error parsing blog template: %v", err)
		return
	}

	err = temp.Execute(w, md)
	if err != nil {
		fmt.Printf("error executing template data: %v", err)
		return
	}
}
