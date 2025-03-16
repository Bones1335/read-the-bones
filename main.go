package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "./static"
	const port = "8080"

	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/static", http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/static/", fsHandler)

	mux.HandleFunc("/", handlerGetIndex)
	mux.HandleFunc("/blog", handlerBlog)
	mux.HandleFunc("/blog/{postTitle}", handlerGetPost)
	mux.HandleFunc("/projects", handlerProjects)
	mux.HandleFunc("/about", handlerAbout)
	mux.HandleFunc("/contact", handlerContact)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
