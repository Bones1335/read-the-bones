package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func main() {
	const filepathRoot = "./static"
	const port = "8080"

	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/static", http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/static/", fsHandler)

	mux.HandleFunc("GET /app/", handlerGetIndex)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

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
