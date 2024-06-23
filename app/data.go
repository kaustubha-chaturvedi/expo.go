package app

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetAllData(srcDir, destDir, templateFile, pagePrefix string, itemsPerPage int,ThemeDir string) ([]interface{}, int) {
	var allData []interface{}

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.Contains(path, ".md") {
			return err
		}

		item := getData(path)
		allData = append(allData, item)
		return nil
	})

	totalPages := (len(allData) + itemsPerPage - 1) / itemsPerPage
	if err != nil {
		log.Fatal(err)
	}

	RenderPages(srcDir, destDir, templateFile, pagePrefix, allData, totalPages, itemsPerPage, ThemeDir)

	return allData, totalPages
}

func getData(path string) interface{} {
	isProject := strings.Contains(path, "projects")
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var content string
	var inHeader bool
	var header, body []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			inHeader = !inHeader
			continue
		}
		if inHeader {
			header = append(header, line)
		} else {
			body = append(body, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if isProject {
		project := parseProject(header)
		return project
	}

	post := parsePost(header)
	post.Slug = strings.ReplaceAll(strings.Split(filepath.Base(path), ".md")[0], " ", "-")
	content = strings.TrimSpace(strings.Join(body, "\n"))
	post.Content = template.HTML(ParseMarkdown(content))
	
	averageWordsPerMinute := 200 
    wordsCount := countWords(content)
    readTimeMinutes := wordsCount / averageWordsPerMinute
    post.ReadTime = fmt.Sprintf("%d min read", readTimeMinutes)
	return post
}

func GetAboutContent(filePath string,ThemeDir string, PublicDir string, Projects []interface{}) (AboutData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return AboutData{}, err
	}
	defer file.Close()
	aboutData := parseAbout(file, Projects)
	RenderPage(filepath.Join(ThemeDir, "about.html"), filepath.Join(PublicDir, "index.html"), aboutData)
	return aboutData, nil
}