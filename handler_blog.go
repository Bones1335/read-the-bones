package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gomarkdown/markdown"
)

func mdToHTML(mdFilePath string) string {
	mdFile, err := os.ReadFile(mdFilePath)
	if err != nil {
		fmt.Printf("error reading markdown file: %v", err)
		return ""
	}

	html := markdown.ToHTML(mdFile, nil, nil)

	return string(html)
}

func handlerBlog(w http.ResponseWriter, r *http.Request) {
	md := mdToHTML("content/connemara_Mathieu_Nicolas.md")

	temp, err := template.ParseFiles("templates/layout.html", "templates/blog.html")
	if err != nil {
		fmt.Printf("error parsing blog template: %v", err)
		return
	}

	err = temp.Execute(w, md)
	if err != nil {
		fmt.Printf("error executing template data: %v\n", err)
		return
	}
}

func handlerGetPost(w http.ResponseWriter, r *http.Request) {
	postTitle := r.PathValue("postTitle")

	postDirectory, err := findPostDirectory(postTitle)
	if err != nil {
		fmt.Printf("couldn't find post: %v\n", err)
		return
	}

	md := mdToHTML(postDirectory + "/index.md")

	temp, err := template.ParseFiles("templates/layout.html", "templates/posts.html")
	if err != nil {
		fmt.Printf("error parsing posts template: %v\n", err)
	}

	err = temp.Execute(w, md)
	if err != nil {
		fmt.Printf("error executing template data: %v\n", err)
		return
	}

}

func findPostDirectory(postTitle string) (string, error) {
	contentDir := "content"
	var postDir string

	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && filepath.Base(path) == postTitle {
			postDir = path
			return fs.SkipAll
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if postDir == "" {
		return "", fmt.Errorf("post not found")
	}

	return postDir, nil
}
