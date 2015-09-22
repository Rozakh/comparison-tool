package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/rozakh/comparator"
	"github.com/rozakh/splunk"
)

const (
	configFileName = "config.json"
	reportFileName = "report.html"
)

var (
	config       configuration
	splunkClient *splunk.Client
	report       *template.Template
)

type configuration struct {
	SplunkHost, Username, Password, CompareHostA, CompareHostB, URLPatternA, URLPatternB string
	ResultsNumber                                                                        int
	CompareElements                                                                      []string
	Groups                                                                               map[string]string
}

type compareResult struct {
	AURL, BURL string
	Diffs      []comparator.Diff
}

func init() {
	report = template.Must(template.New("template").Parse(templateText))
}

func main() {
	var configFile []byte
	workDir, err := os.Getwd()
	checkError(err)
	_, err = os.Stat(workDir + "/" + configFileName)
	if err != nil {
		jsonConfig, err := json.MarshalIndent(config, "", "    ")
		checkError(err)
		ioutil.WriteFile(workDir+"/"+configFileName, jsonConfig, 0660)
	}
	configFile, err = ioutil.ReadFile(workDir + "/" + configFileName)
	checkError(err)
	json.Unmarshal(configFile, &config)
	results := compareResponses()
	generateReport(results)
}

func compareResponses() []compareResult {
	var compareResults []compareResult
	urlsToCompare := make(map[string]string)
	splunkClient = splunk.New(config.SplunkHost)
	err := splunkClient.Login(config.Username, config.Password)
	checkError(err)
	searchQuery := getSearchQueryFromTemplate(config.URLPatternA)
	results, err := splunkClient.Search(searchQuery, config.ResultsNumber)
	checkError(err)
	for _, result := range results {
		urlA := replacePlaceholdersWithValues(config.URLPatternA, result)
		urlB := replacePlaceholdersWithValues(config.URLPatternB, result)
		urlsToCompare[config.CompareHostA+urlA] = config.CompareHostB + urlB
	}
	for k, v := range urlsToCompare {
		result, err := comparator.Compare(k, v, config.CompareElements)
		checkError(err)
		if len(result) != 0 {
			compareResults = append(compareResults, compareResult{k, v, result})
		}
	}
	return compareResults
}

func generateReport(results []compareResult) {
	workDir, err := os.Getwd()
	checkError(err)
	file, err := os.OpenFile(workDir+"/"+reportFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0660)
	checkError(err)
	defer file.Close()
	report.Execute(file, results)
}

func getSearchQueryFromTemplate(template string) string {
	template = strings.Replace(template, "?", "\\?", -1)
	for strings.Contains(template, "{{") {
		group, groupName := getGroupAndName(template)
		replacement := config.Groups[groupName]
		if replacement != "" {
			replacement = fmt.Sprintf("(?<%s>%s?)", groupName, replacement)
		} else {
			replacement = fmt.Sprintf("(?<%s>.*?)", groupName)
		}
		template = strings.Replace(template, group, replacement, -1)
	}
	return fmt.Sprintf("GET (?<prefix>.*?)%s(?<suffix>.*?) HTTP", template)
}

func replacePlaceholdersWithValues(template string, values map[string]string) string {
	for strings.Contains(template, "{{") {
		group, groupName := getGroupAndName(template)
		template = strings.Replace(template, group, values[groupName], -1)
	}
	return values["prefix"] + template + values["suffix"]
}
func getGroupAndName(template string) (string, string) {
	groupStart := strings.Index(template, "{{")
	groupEnd := strings.Index(template, "}}")
	groupName := template[groupStart+2 : groupEnd]
	group := fmt.Sprintf("{{%s}}", groupName)
	return group, groupName
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
