package main

import (
	"bufio"
	"fmt"
	"github.com/reiver/go-porterstemmer"
	"io/ioutil"
	"os"
	"strings"
)

var stopWords string

func isStopWord(word string) bool {
	return strings.Contains(stopWords, strings.ToLower(word))
}

func stem(word string) string {
	return porterstemmer.StemString(word)
}

func customPostTokenizerClean(word string) string {
	// remove end punctuations
	word = strings.TrimRight(word, ",.;:?!\")]}(-'")
	word = strings.TrimLeft(word, "([{\"'")
	return word
}

func customPreTokenizerClean(input string) string {
	for _, s := range []string{",", ";", ":", "=", "--",
		"__", "**", "##", "&lt", "&gt", "(",
		"\u2018",   // single left quote
		"\u2019",   // single right quote
		"\u201C",   // double left quote
		"\u201D"} { // double right quote
		input = strings.Replace(input, s, " ", -1)
	}
	return filterLogLines(input)
}

func isIntegerNumber(word string) bool {
	for _, c := range word {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func filter(input string) []string {
	input = customPreTokenizerClean(input)
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords) // need a proper split to clean out punctuations

	var words []string
	for scanner.Scan() {
		token := customPostTokenizerClean(scanner.Text())
		if !isStopWord(token) && !isIntegerNumber(token) {
			words = append(words, stem(token))
		}
	}
	return words
}

func save(filename, content string) {
	fmt.Println("Saving to", filename)
	ioutil.WriteFile(filename, []byte(content), 0666)
}

func getFilename(filename string) string {
	return filename[:len(filename)-len(".txt")] + ".cln"
}

func filterLogLines(content string) string {
	lines := strings.Split(content, "\n")
	result := make([]string, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(trimmed, "at org") ||
			strings.HasPrefix(trimmed, "at sun") ||
			strings.HasPrefix(trimmed, "at java") {
			continue
		}
		result = append(result, trimmed)
	}
	return strings.Join(result, "\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide input text file.")
		os.Exit(1)
	}

	buffer, _ := ioutil.ReadFile("common-english-words.txt")
	stopWords = string(buffer)
	filename := os.Args[1]
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error loading file", filename)
		os.Exit(1)
	}
	words := filter(string(content))
	save(getFilename(filename), strings.Join(words, " "))
}
