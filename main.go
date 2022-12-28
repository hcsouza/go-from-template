package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/Songmu/prompter"
)

type ProjectGenerator struct {
	PROJECT_PATH        string
	PROJECT_NAME        string
	PROJECT_DESTINATION string
	TEMPLATE_NAME       string
}

func main() {

	templateName := (&prompter.Prompter{
		Choices:    []string{"go-rest-api-template"},
		Default:    "go-rest-api-template",
		Message:    "Choose project template",
		IgnoreCase: true,
	}).Prompt()

	projectPath := (&prompter.Prompter{
		Default:    "github.com/user",
		Message:    "Choose repository path",
		IgnoreCase: true,
	}).Prompt()

	projectName := (&prompter.Prompter{
		Default:    "generated",
		Message:    "Choose generated project name",
		IgnoreCase: true,
	}).Prompt()

	projectDestination := (&prompter.Prompter{
		Default:    "/tmp",
		Message:    "Choose destination path",
		IgnoreCase: true,
	}).Prompt()

	generated := ProjectGenerator{projectPath, projectName, projectDestination, templateName}
	generateProject(generated)
}

func generateProject(generated ProjectGenerator) {
	files, err := searchAllFiles(fmt.Sprintf("./%s", generated.TEMPLATE_NAME))
	if err != nil {
		fmt.Println("RunSearchError: ", err)
	}
	for _, fileName := range files {
		RunFile(fileName, generated)
	}
}

func RunFile(templatePath string, generated ProjectGenerator) {

	tmplContainer, err := template.New("generator").ParseFiles(templatePath)
	if err != nil {
		fmt.Println("create tmp Error: ", err)
	}

	rootPath := strings.Replace(templatePath, generated.TEMPLATE_NAME, generated.PROJECT_NAME, 1)
	fileName := strings.Replace(rootPath, ".tmpl", "", 1)

	if len(generated.PROJECT_DESTINATION) > 0 {
		fileName = fmt.Sprintf("%s/%s", generated.PROJECT_DESTINATION, fileName)
	}

	dir, _ := filepath.Split(fileName)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
	}

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}

	fmt.Println("Generating file: ", fileName)
	err = tmplContainer.Execute(f, generated)
	if err != nil {
		panic(err)
	}
	f.Close()
}

func searchAllFiles(searchDir string) ([]string, error) {
	fileList := make([]string, 0)
	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() == false && strings.Contains(path, ".tmpl") == true {
			fileList = append(fileList, path)
		}
		return err

	})

	if e != nil {
		panic(e)
	}
	return fileList, nil
}
