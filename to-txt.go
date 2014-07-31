package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func save(filename, content string) {
	fmt.Println("saving: " + filename)
	ioutil.WriteFile(filename, []byte(content), 0666)
}

func getFilename(filename string) string {
	return filename[:len(filename)-len(".xml")] + ".txt"
}

func convert(content string) (string, error) {
	v := new(RSS)
	err := xml.Unmarshal([]byte(content), &v)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", v.RssChannel.Issue), nil
}

type RSS struct {
	XMLName    xml.Name `xml:"rss"`
	RssChannel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name
	Issue   Item `xml:"item"`
}

type Item struct {
	Description string  `xml:"description"`
	Key         string  `xml:"key"`
	Summary     string  `xml:"summary"`
	Component   string  `xml:"component"`
	Comments    Comment `xml:"comments"`
}

func (i Item) String() string {
	return fmt.Sprintf("%s - %s\nComponent: %s\n\n%s\n\n%v",
		i.Key, clean(i.Summary), i.Component, clean(i.Description), i.Comments)
}

type Comment struct {
	Content []string `xml:"comment"`
}

func (c Comment) String() string {
	comments := make([]string, len(c.Content), len(c.Content))
	for i, comment := range c.Content {
		comments[i] = clean(comment)
	}
	return fmt.Sprintf("\n%s", strings.Join(comments, "\n"))
}

func clean(s string) string {
	// todo use html parser to clean
	re := regexp.MustCompile("</?[^>]+>")
	return re.ReplaceAllString(s, "")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide input file name")
		os.Exit(1)
	}

	filename := os.Args[1]
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error loading file", filename)
		os.Exit(1)
	}
	converted, err := convert(string(content))
	if err != nil {
		fmt.Println("Error converting file", filename)
		os.Exit(1)
	}
	save(getFilename(filename), converted)
}
