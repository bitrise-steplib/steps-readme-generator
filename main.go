package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/templateutil"
	"github.com/bitrise-io/stepman/models"
	"gopkg.in/yaml.v2"
)

type config struct {
	ExampleSection string `env:"example_section"`
	ContribSection string `env:"contrib_section"`
}

//go:embed README.md.gotemplate
var readmeTemplate string

type templateInventory struct {
	Step           models.StepModel
	ExampleSection string
	ContribSection string
}

func createBackup() error {
	err := os.Rename("README.md", "README.md.backup")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to rename README.md to README.md.backup: %w", err)
	}
	log.Donef("Created backup as README.md.backup")

	return nil
}

func parseStep() (models.StepModel, error) {
	fileContents, err := ioutil.ReadFile("step.yml")
	if err != nil {
		return models.StepModel{}, fmt.Errorf("failed to open step.yml: %w", err)
	}

	stepConfig := models.StepModel{}
	if err = yaml.Unmarshal(fileContents, &stepConfig); err != nil {
		return models.StepModel{}, fmt.Errorf("failed to parse step.yml: %w", err)
	}

	return stepConfig, nil
}

func readSections(stepConfig config) (exampleSection, contribSection string, err error) {
	readSection := func(name, pth string) (string, error) {
		log.Infof("Using %s section from %s", name, pth)
		cont, err := ioutil.ReadFile(pth)
		if err != nil {
			return "", err
		}
		return string(cont), nil
	}
	if stepConfig.ExampleSection != "" {
		exampleSection, err = readSection("example", stepConfig.ExampleSection)
	}

	if err == nil && stepConfig.ContribSection != "" {
		contribSection, err = readSection("contrib", stepConfig.ContribSection)
	}

	return exampleSection, contribSection, nil
}

func markdownTableCompatibleString(text string) string {
	withoutNewlines := strings.ReplaceAll(text, "\n", " ")
	escapedPipes := strings.ReplaceAll(withoutNewlines, "|", "\\|")
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

func githubName(repoURL string) string {
	return strings.Split(repoURL, "github.com/")[1]
}

// hasDefault tells if an input has a default value or not.
// Unlike using {{ if $default }}, this handles cases when a bool input is false or when an int input is 0
func hasDefault(defaultInputValue interface{}) bool {
	if defaultInputValue == "" {
		return false
	}
	return defaultInputValue != nil
}

func renderTemplate(step models.StepModel, exampleSection, contribSection string) (string, error) {
	funcMap := template.FuncMap{
		"markdownTableCompatibleString": markdownTableCompatibleString,
		"flagList":                      flagList,
		"githubName":                    githubName,
		"hasDefault":                    hasDefault,
	}

	inventory := templateInventory{
		Step:           step,
		ExampleSection: exampleSection,
		ContribSection: contribSection,
	}

	readmeContent, err := templateutil.EvaluateTemplateStringToString(readmeTemplate, inventory, funcMap)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate template: %w", err)
	}

	return readmeContent, nil
}

func writeReadme(contents string) error {
	if err := ioutil.WriteFile("README.md", []byte(contents), 0644); err != nil {
		return fmt.Errorf("failed to write README contents to file: %w", err)
	}
	return nil
}

func mainR() error {
	var stepConfig config
	if err := stepconf.Parse(&stepConfig); err != nil {
		return err
	}
	stepconf.Print(stepConfig)
	fmt.Println()

	log.Infof("Generating README.md from step.yml data")

	if err := createBackup(); err != nil {
		return err
	}

	stepData, err := parseStep()
	if err != nil {
		return err
	}

	exampleSection, contribSection, err := readSections(stepConfig)
	if err != nil {
		return err
	}

	readmeContents, err := renderTemplate(stepData, exampleSection, contribSection)
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
