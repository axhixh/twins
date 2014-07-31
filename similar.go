package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
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

func save(filename string, content map[string]float64) {
	fmt.Println("Saving file", filename)
	outFile, err := os.Create(filename)
	defer outFile.Close()

	if err != nil {
		fmt.Println("Unable to save file", filename)
		os.Exit(1)
	}

	for k, v := range content {
		outFile.WriteString(fmt.Sprintf("%.4f %s\n", v, k))
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Please provide the directory containing the files")
		os.Exit(1)
	}
	fmt.Println("Loading files")
	files := loadFiles(os.Args[1])
	result := make(map[string]float64)
	var filenames []string
	for k, _ := range files {
		filenames = append(filenames, k)
	}

	threshold := 0.8

	fmt.Println("Calculating similarity")
	for i, f1 := range filenames {
		for j := i + 1; j < len(filenames); j++ {
			value := similarity(files[f1], files[filenames[j]])
			if value > threshold {
				result[fmt.Sprintf("%s:%s", f1, filenames[j])] = value
			}
		}
	}
	save("similarity.txt", result)
}

func loadFiles(path string) map[string]map[string]float64 {
	files := make(map[string]map[string]float64)
	loader := func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".idf") {
			content, err := load(path)
			if err != nil {
				return err
			}
			files[info.Name()] = content
		}
		return nil
	}
	err := filepath.Walk(path, loader)
	if err != nil {
		fmt.Println("Error loading idf files")
		os.Exit(1)
	}
	return files
}

func similarity(doc1, doc2 map[string]float64) float64 {
	allWords := words(doc1, doc2)
	vec1 := toVector(allWords, doc1)
	vec2 := toVector(allWords, doc2)
	return cosineSimilarity(vec1, vec2)
}

func words(doc1, doc2 map[string]float64) []string {
	keys := make(map[string]struct{})
	for k, _ := range doc1 {
		keys[k] = struct{}{}
	}
	for k, _ := range doc2 {
		keys[k] = struct{}{}
	}
	var result []string
	for k, _ := range keys {
		result = append(result, k)
	}
	return result
}

func toVector(words []string, doc map[string]float64) []float64 {
	result := make([]float64, len(words))
	for i, v := range words {
		result[i] = doc[v]
	}
	return result
}

func cosineSimilarity(vec1, vec2 []float64) float64 {
	return dotProduct(vec1, vec2) / (norm(vec1) * norm(vec2))
}

func dotProduct(vec1, vec2 []float64) float64 {
	var sum float64
	sum = 0
	for i := range vec1 {
		sum += (vec1[i] * vec2[i])
	}
	return sum
}

func norm(vec []float64) float64 {
	var sum float64
	sum = 0
	for _, v := range vec {
		sum += (v * v)
	}
	return math.Sqrt(sum)
}
