package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func handlerAbout(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/layout.html", "templates/about.html")
	if err != nil {
		fmt.Printf("error parsing about template: %v", err)
		return
	}

	err = temp.Execute(w, "")
	if err != nil {
		fmt.Printf("error executing about template data: %v", err)
		return
	}
}
