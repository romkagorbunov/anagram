package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
)

func uniqueSortedStrings(input []string) []string {
	// Sort input first
	sort.Strings(input)

	// Result slice
	var result []string

	for i, s := range input {
		// Skip duplicates
		if i == 0 || s != input[i-1] {
			result = append(result, s)
		}
	}

	return result
}

func main() {
	filePath := "voyna_i_mir.txt"

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	content := string(data)

	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	words := []string{}

	for _, line := range lines {
		curw := ""
		for _, c := range line {
			if unicode.IsLetter(c) {
				curw = curw + string(c)
			} else {
				if curw != "" {
					words = append(words, strings.ToLower(curw))
					curw = ""
				}
			}
		}
	}
	words = uniqueSortedStrings(words)
	outFilePath := "words.txt"
	outstring := ""
	for _, word := range words {
		outstring = outstring + "\n" + word
	}

	err = ioutil.WriteFile(outFilePath, []byte(outstring), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
