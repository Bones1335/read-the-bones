package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"text/template"
)

func handlerGetIndex(w http.ResponseWriter, r *http.Request) {
	directories, err := getDirectoryNames("content/posts")
	if err != nil {
		fmt.Printf("error getting directory names: %v", err)
		return
	}

	var posts []PostMetaData

	for _, dir := range directories {
		path := fmt.Sprintf("%v/index.md", dir)
		_, err := os.Open(path)
		if os.IsNotExist(err) {
			continue
		}

		metaData, err := parseMarkdown(path)
		if err != nil {
			fmt.Printf("error parsing markdown data: %v", err)
			return
		}
		posts = append(posts, metaData)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	temp, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	if err != nil {
		fmt.Printf("error parsing html templates: %v", err)
		return
	}

	err = temp.Execute(w, posts)
	if err != nil {
		fmt.Printf("problem executing template data: %v", err)
		return
	}
}
