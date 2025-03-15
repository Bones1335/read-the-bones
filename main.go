package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gomarkdown/markdown"
)

func main() {
	const filepathRoot = "./static"
	const port = "8080"

	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/static", http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/static/", fsHandler)

	mux.HandleFunc("/", handlerGetIndex)
	mux.HandleFunc("/blog", handlerBlog)
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

func mdToHTML() string {
	mdFile, err := os.ReadFile("content/connemara_Mathieu_Nicolas.md")
	if err != nil {
		fmt.Printf("error reading markdown file: %v", err)
		return ""
	}

	html := markdown.ToHTML(mdFile, nil, nil)

	return string(html)
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
