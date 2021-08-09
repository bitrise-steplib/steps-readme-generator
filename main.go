package main

import (
	_ "embed"
	"fmt"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/templateutil"
	"github.com/bitrise-io/stepman/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

//go:embed README.md.gotemplate
var readmeTemplate string

func createBackup() error {
	err := os.Rename("README.md", "README.md.backup")
	if err != nil {
		return fmt.Errorf("failed to rename README.md to README.md.backup, error: %s", err)
	}
	return nil
}

func parseStep() (models.StepModel, error) {
	fileContents, err := ioutil.ReadFile("step.yml")
	if err != nil {
		return models.StepModel{}, fmt.Errorf("failed to open step.yml, error: %s", err)
	}

	stepConfig := models.StepModel{}
	if err = yaml.Unmarshal(fileContents, &stepConfig); err != nil {
		return models.StepModel{}, fmt.Errorf("failed to parse step.yml, error: %s", err)
	}

	return stepConfig, nil
}

// clean makes Markdown table–compatible text from the input by replacing line breaks with spaces and escaping special
// characters.
func clean(text string) string {
	withoutNewlines := strings.Replace(text, "\n", " ", -1)
	escapedPipes := strings.Replace(withoutNewlines, "|", "\\|", -1)
	return escapedPipes
}

func flagList(isRequired, isSensitive interface{}) string {
	var flags []string
	if isRequired == true {
		flags = append(flags, "required")
	}
	if isSensitive == true {
		flags = append(flags, "sensitive")
	}

	return strings.Join(flags, ", ")
}

func renderTemplate(stepConfig models.StepModel) (string, error) {
	funcMap := template.FuncMap{
		"clean":    clean,
		"flagList": flagList,
	}
	readmeContent, err := templateutil.EvaluateTemplateStringToString(readmeTemplate, stepConfig, funcMap)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate template, error: %s", err)
	}

	return readmeContent, nil
}

func writeReadme(contents string) error {
	if err := ioutil.WriteFile("README.md", []byte(contents), 0644); err != nil {
		return fmt.Errorf("failed to write README contents to file, error: %s", err)
	}
	return nil
}

func mainR() error {
	log.Infof("Generating README.md from step.yml data")

	if err := createBackup(); err != nil {
		return err
	}
	log.Donef("Created backup as README.md.backup")

	stepConfig, err := parseStep()
	if err != nil {
		return err
	}

	readmeContents, err := renderTemplate(stepConfig)
	if err != nil {
		return err
	}

	err = writeReadme(readmeContents)
	if err != nil {
		return err
	}
	log.Donef("README.md generated successfully")

	return nil
}

func main() {
	if err := mainR(); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
}
