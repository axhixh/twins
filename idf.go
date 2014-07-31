package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func load(filename string) (map[string]float64, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	result := make(map[string]float64)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Error splitting line %s", line)
		}
		value, _ := strconv.ParseFloat(parts[1], 64)
		result[parts[0]] = value
	}
	return result, nil
}

func calculate(tf, idf map[string]float64) map[string]float64 {
	result := make(map[string]float64, len(tf))
	for word, value := range tf {
		result[word] = value * idf[word]
	}
	return result
}

func save(filename string, content map[string]float64) {
	fmt.Println("Saving file", filename)
	outFile, err := os.Create(filename)
	defer outFile.Close()

	if err != nil {
		fmt.Println("Unable to save file", filename)
		os.Exit(1)
	}

	for k, v := range content {
		outFile.WriteString(fmt.Sprintf("%s %v\n", k, v))
	}
}

func outputFilename(filename string) string {
	return filename[:len(filename)-len(".tf")] + ".idf"
}

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "input", "", "File containg term frequency")
	var idfFile string
	flag.StringVar(&idfFile, "idf", "corpus.idf", "File containing inverse document frequency")

	flag.Parse()

	tf, err := load(inputFile)
	if err != nil {
		fmt.Println("Unable to load", inputFile, err)
		os.Exit(1)
	}
	idf, err := load(idfFile)
	if err != nil {
		fmt.Println("Unable to load", idfFile, err)
		os.Exit(1)
	}

	result := calculate(tf, idf)
	save(outputFilename(inputFile), result)
}
