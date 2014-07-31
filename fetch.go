package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

func getContent(url, username, password string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

type Params struct {
	dataFolder string
	username   string
	password   string
	issueId    string
	baseUrl    string
}

func getParams() (*Params, error) {
	var output string
	flag.StringVar(&output, "output", "data", "Output folder for fetched issue")
	var username string
	flag.StringVar(&username, "user", "jira-user", "Username for Jira")
	var password string
	flag.StringVar(&password, "password", "", "Password for Jira")
	var issue string
	flag.StringVar(&issue, "issue", "", "Jira Issue Id")
	var baseUrl string
	flag.StringVar(&baseUrl, "url", "http://jiraserver.com", "URL to the Jira server")
	flag.Parse()

	if len(password) == 0 {
		return new(Params), fmt.Errorf("The password cannot be empty.")
	}

	issue = strings.ToUpper(strings.Trim(issue, " "))
	if len(issue) == 0 {
		return new(Params), fmt.Errorf("The issue cannot be empty.")
	}

	return &Params{output, username, password, issue, baseUrl}, nil
}

func getUrl(baseUrl, issue string) string {
	return fmt.Sprintf("%s/si/jira.issueviews:issue-xml/%s/%s.xml",
		baseUrl, issue, issue)
}

func save(filename string, content []byte) {
	fmt.Printf("Saving to file %s\n", filename)
	ioutil.WriteFile(filename, content, 0666)
}

func main() {
	params, err := getParams()
	if err != nil {
		panic(err)
	}
	url := getUrl(params.baseUrl, params.issueId)
	content, err := getContent(url, params.username, params.password)
	if err != nil {
		panic(err)
	}
	filename := path.Join(params.dataFolder, params.issueId+".xml")
	save(filename, content)
}
