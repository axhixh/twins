package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func count(tdf map[string]int, filename string) error {
	// fixme: process the content of the file using a stream/Reader
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		words := strings.Split(line, " ")
		tdf[words[0]] = tdf[words[0]] + 1
	}
	return nil
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
		fmt.Println("Please provide the directory containg *.tf files")
		os.Exit(1)
	}

	noOfDocs := 0
	folder := os.Args[1]
	tdf := make(map[string]int)
	counter := func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".tf") {
			if err := count(tdf, path); err != nil {
				return err
			}
			noOfDocs += 1
		}
		return nil
	}

	err := filepath.Walk(folder, counter)
	if err != nil {
		fmt.Println("error generating idf")
		os.Exit(1)
	}

	idf := make(map[string]float64, len(tdf))
	for word, count := range tdf {
		idf[word] = math.Log(float64(noOfDocs) / float64(count))
	}

	save("corpus.idf", idf)
}
