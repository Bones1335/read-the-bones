package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func handlerGetIndex(w http.ResponseWriter, r *http.Request) {
	text := "Hello internet! Welcome to Read the Bones :)"

	temp, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		fmt.Printf("error parsing html templates: %v", err)
		return
	}

	err = temp.Execute(w, text)
	if err != nil {
		fmt.Printf("problem executing template data: %v", err)
		return
	}
}
