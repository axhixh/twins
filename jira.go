package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Jira struct {
	BaseURL  string
	Username string
	Password string
}

func (j Jira) GetIssue(issue string) ([]byte, error) {
	url := j.getUrl(issue)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(j.Username, j.Password)
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

func (j Jira) getUrl(issue string) string {
	return fmt.Sprintf("%s/si/jira.issueviews:issue-xml/%s/%s.xml", j.BaseURL, issue, issue)
}
