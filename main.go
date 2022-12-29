package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Songmu/prompter"
)

type ProjectGenerator struct {
	PROJECT_PATH        string
	PROJECT_NAME        string
	PROJECT_DESTINATION string
	TEMPLATE_NAME       string
	TEMPLATE_GIT_PATH   string
}

var DEFAULT_TEMPLATES map[string]string

func init() {
	DEFAULT_TEMPLATES = make(map[string]string)
	DEFAULT_TEMPLATES["go-rest-api-template"] = "git@github.com:hcsouza/go-rest-api-template.git"
}

func main() {

	templateName := (&prompter.Prompter{
		//Choices:    []string{"go-rest-api-template"},
		Default:    "go-rest-api-template",
		Message:    "Choose a default project template or an external",
		IgnoreCase: true,
	}).Prompt()

	templateGitPath, ok := DEFAULT_TEMPLATES[templateName]
	if !ok {
		templateGitPath = (&prompter.Prompter{
			Message:    "Pass the template git path: git@github.com:user/user.git",
			IgnoreCase: true,
		}).Prompt()
	}

	projectPath := (&prompter.Prompter{
		Default:    "github.com/user",
		Message:    "Choose destination repository path",
		IgnoreCase: true,
	}).Prompt()

	projectName := (&prompter.Prompter{
		Default:    "new-go-project",
		Message:    "Choose generated project name",
		IgnoreCase: true,
	}).Prompt()

	dirDestination, _ := os.Getwd()
	projectDestination := (&prompter.Prompter{
		Default:    dirDestination,
		Message:    "Choose destination path",
		IgnoreCase: true,
	}).Prompt()

	generated := ProjectGenerator{projectPath, projectName, projectDestination, templateName, templateGitPath}
	cloneTemplate(generated)
	generateProject(generated)
}

func generateProject(generated ProjectGenerator) {
	files, err := searchAllFiles(fmt.Sprintf("./Templates/%s", generated.TEMPLATE_NAME))
	if err != nil {
		fmt.Println("RunSearchError: ", err)
	}
	for _, fileName := range files {
		RunFile(fileName, generated)
	}
}

func cloneTemplate(generated ProjectGenerator) {
	fmt.Println("Cloning the template!!")
	fmt.Println()
	dirname, _ := os.Getwd()
	clonePath := fmt.Sprintf("%s/Templates/%s", dirname, generated.TEMPLATE_NAME)
	if _, err := os.Stat(clonePath); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", generated.TEMPLATE_GIT_PATH, clonePath)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func RunFile(templatePath string, generated ProjectGenerator) {

	tmplContainer, err := template.New("generator").ParseFiles(templatePath)
	if err != nil {
		fmt.Println("create tmp Error: ", err)
	}

	rootPath := strings.Replace(templatePath, generated.TEMPLATE_NAME, generated.PROJECT_NAME, 1)
	rootPath = strings.Replace(rootPath, "Templates/", "", 1)
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
		if !f.IsDir() && strings.Contains(path, ".tmpl") {
			fileList = append(fileList, path)
		}
		return err

	})

	if e != nil {
		panic(e)
	}
	return fileList, nil
}
