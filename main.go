package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

const firstLetter = 'а'
const once = 1000 * 1000
const numLetters = 'я' - firstLetter + 1

// Returns true if the next permutation was found; false if it's the last permutation

func nextUniquePermutation(s string) (string, bool) {
	runes := []rune(s)
	n := len(runes)

	// Step 1: Find pivot — the largest index i such that runes[i] < runes[i+1]
	i := n - 2
	for i >= 0 && runes[i] >= runes[i+1] {
		i--
	}
	if i < 0 {
		// Already the last permutation (in lexicographic order). Reset to first.
		reverseRunes(runes, 0, n-1)
		return string(runes), false
	}

	// Step 2: Find the rightmost successor to pivot in suffix ( > runes[i] )
	j := n - 1
	for j > i && runes[j] <= runes[i] {
		j--
	}

	if i == j {
		fmt.Println("Fuck")
	}

	// Swap pivot and successor
	runes[i], runes[j] = runes[j], runes[i]

	// Step 3: Reverse the suffix (i+1 ... end)
	reverseRunes(runes, i+1, n-1)

	// Now we have a new permutation; but it might equal the old string if duplicates
	newS := string(runes)
	return newS, true
}

func reverseRunes(runes []rune, left, right int) {
	for left < right {
		runes[left], runes[right] = runes[right], runes[left]
		left++
		right--
	}
}

func sortString(s string) string {
	runes := []rune(s)

	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})

	return string(runes)
}

type Node struct {
	move       []*Node
	isTerminal bool
}

func getId(c rune) rune {
	return c - firstLetter
}

func NewNode() *Node {
	n := Node{
		make([]*Node, numLetters),
		false,
	}
	return &n
}

func addWord(s string, root *Node) {
	cur := root
	for _, c := range s {
		if cur.move[getId(c)] == nil {
			cur.move[getId(c)] = NewNode()
		}
		cur = cur.move[getId(c)]
	}
	cur.isTerminal = true
}

func tryFind(name string, root *Node) bool {
	cur := root
	for _, c := range name {
		if c == ' ' {
			if cur.isTerminal {
				cur = root
			} else {
				return false
			}
		} else if cur.move[getId(c)] == nil {
			return false
		} else {
			cur = cur.move[getId(c)]
		}
	}
	return cur.isTerminal
}

func dfs(node *Node, have []int, spaces int, cur string, root *Node, writer *bufio.Writer) {
	// fmt.Println("wtf", cur)
	if node == nil {
		return
	}
	if spaces > 0 && node.isTerminal {
		nhave := make([]int, len(have)) // create a new slice with the same length
		copy(nhave, have)
		dfs(root, nhave, spaces-1, cur+" ", root, writer)
	}
	ss := 0
	for i := rune(0); i < numLetters; i++ {
		ss += have[i]
		if have[i] != 0 && node.move[i] != nil {
			nhave := make([]int, len(have)) // create a new slice with the same length
			copy(nhave, have)
			nhave[i]--
			ncur := cur + string(rune(i+firstLetter))
			dfs(node.move[i], nhave, spaces, ncur, root, writer)
		}
	}
	if ss == 0 && node.isTerminal {
		_, err := writer.WriteString(cur + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func genOutputName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

func main() {
	fmt.Println("Some")
	filePath := "utf-8.txt"

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Convert []byte to string (Go strings are UTF-8 by default)
	content := string(data)

	root := NewNode()
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	for _, line := range lines {
		line = strings.ToLower(line)
		goodWord := true
		for _, c := range line {
			if getId(c) < 0 || getId(c) >= numLetters {
				goodWord = false
			}
		}
		if !goodWord {
			continue
		}
		fmt.Println("Adding ", line)
		addWord(line, root)
	}

	namesData, err := os.ReadFile("names")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	names := string(namesData)
	nameLines := strings.Split(strings.ReplaceAll(names, "\r\n", "\n"), "\n")

	var wg = sync.WaitGroup{}

	for _, name := range nameLines {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			cntSpaces := 0
			for _, c := range name {
				if c == ' ' {
					cntSpaces++
				}
			}
			outputName := genOutputName(name)

			name = strings.ToLower(strings.ReplaceAll(name, " ", ""))
			fmt.Println(name)

			file, err := os.Create(outputName)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			writer := bufio.NewWriter(file)

			name = sortString(name)
			have := make([]int, numLetters)
			for _, c := range name {
				have[getId(c)]++
			}
			dfs(root, have, cntSpaces, "", root, writer)

			err = writer.Flush()
			if err != nil {
				log.Fatal(err)
			}
		}(name)
	}
	wg.Wait()
}
