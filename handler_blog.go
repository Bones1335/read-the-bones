package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
	extensions := parser.CommonExtensions | parser.HardLineBreak | parser.Footnotes
	p := parser.NewWithExtensions(extensions)

	html := markdown.ToHTML([]byte(mdFile), p, nil)

	return string(html)
}

func handlerBlog(w http.ResponseWriter, r *http.Request) {
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

	temp, err := template.ParseFiles("templates/layout.html", "templates/blog.html")
	if err != nil {
		fmt.Printf("error parsing blog template: %v", err)
		return
	}

	err = temp.Execute(w, posts)
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

func getDirectoryNames(postsDir string) ([]string, error) {
	var dirMap []string

	err := filepath.WalkDir(postsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			dirMap = append(dirMap, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dirMap, nil

}
