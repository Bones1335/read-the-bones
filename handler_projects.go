package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func handlerProjects(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/layout.html", "templates/projects.html")
	if err != nil {
		fmt.Printf("error parsing projects template: %v", err)
		return
	}

	err = temp.Execute(w, "")
	if err != nil {
		fmt.Printf("error executing projects template data: %v", err)
		return
	}
}
