package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func count(input string) map[string]float64 {
	parts := strings.Fields(input)
	m := make(map[string]int)
	totalWords := 0
	for _, f := range parts {
		m[f] += 1
		totalWords += 1
	}

	f := make(map[string]float64)
	for k, v := range m {
		f[k] = float64(v) / float64(totalWords)
	}
	return f
}

func save(filename string, counts map[string]float64) {
	fmt.Println("Saving file", filename)
	outFile, err := os.Create(filename)
	defer outFile.Close()

	if err != nil {
		fmt.Println("Unable to write file", filename)
		os.Exit(1)
	}

	for k, v := range counts {
		outFile.WriteString(fmt.Sprintf("%s %v\n", k, v))
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide the file")
		os.Exit(1)
	}

	filename := os.Args[1]
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Unable to read file", filename)
		os.Exit(1)
	}

	result := count(string(content))
	filename = filename[:len(filename)-len(".cln")] + ".tf"
	save(filename, result)
}
