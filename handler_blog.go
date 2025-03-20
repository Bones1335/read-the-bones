package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/yaml.v3"
)

type PostMetaData struct {
	Title      string    `yaml:"title"`
	Date       time.Time `yaml:"date"`
	URL        string    `yaml:"url"`
	Categories []string  `yaml:"categories"`
	Tags       []string  `yaml:"tags"`
	Content    string
}

func parseMarkdown(mdFilePath string) (PostMetaData, error) {
	mdFileData, err := os.ReadFile(mdFilePath)
	if err != nil {
		fmt.Printf("error reading markdown file: %v\n", err)
		return PostMetaData{}, err
	}

	contentStr := string(mdFileData)

	re := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`)
	matches := re.FindStringSubmatch(contentStr)

	var metaData PostMetaData
	var markdownBody string

	if len(matches) == 3 {
		err := yaml.Unmarshal([]byte(matches[1]), &metaData)
		if err != nil {
			return PostMetaData{}, fmt.Errorf("error parsing YAML: %w", err)
		}
		markdownBody = matches[2]
	} else {
		markdownBody = contentStr
	}

	metaData.Content = mdToHTML(markdownBody)

	return metaData, nil
}

func mdToHTML(mdFile string) string {
	extensions := parser.CommonExtensions | parser.HardLineBreak
	p := parser.NewWithExtensions(extensions)

	html := markdown.ToHTML([]byte(mdFile), p, nil)

	return string(html)
}

func handlerBlog(w http.ResponseWriter, r *http.Request) {
	metaData, err := parseMarkdown("content/posts/2025/reacher-temporarily-solved-my-rut/index.md")
	if err != nil {
		fmt.Printf("error parsing markdown data: %v", err)
		return
	}

	md := postContainer(metaData)

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

	metaData, err := parseMarkdown(postDirectory + "/index.md")
	if err != nil {
		fmt.Printf("Error parsing Markdown: %v\n", err)
		return
	}

	fmt.Printf("Title: %v\nTags: %v\n", metaData.Title, metaData.Tags)

	temp, err := template.ParseFiles("templates/layout.html", "templates/posts.html")
	if err != nil {
		fmt.Printf("error parsing posts template: %v\n", err)
	}

	err = temp.Execute(w, metaData)
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

func postContainer(metaData PostMetaData) string {
	title := mdToHTML(metaData.Title)
	
	date := mdToHTML(metaData.Date.Format("2006-02-25"))

	containerDiv := fmt.Sprintf("<div class=\"blog-post-card\">%v%v</div>", title, date)

	return containerDiv
}
