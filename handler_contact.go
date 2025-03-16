package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func handlerContact(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/layout.html", "templates/contact.html")
	if err != nil {
		fmt.Printf("error parsing contact template: %v", err)
		return
	}

	err = temp.Execute(w, nil)
	if err != nil {
		fmt.Printf("error executing contact template data: %v", err)
		return
	}
}
