package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var charset string
var N int

var permutations = make(map[string]bool)

var words = make(map[string]bool)
var minLen = 3
var wordsChecked = 0

func swap(l *[]string, index1, index2 int) {
	var t string
	t = (*l)[index1]
	(*l)[index1] = (*l)[index2]
	(*l)[index2] = t
}

func generatePermutations(l []string, n int) {
	if N == n {
		permutations[strings.Join(l, "")] = true
	} else {
		for i := n; i <= N; i += 1 {
			swap(&l, i, n)
			generatePermutations(l, n+1)
			swap(&l, i, n)
		}
	}
}

func generateSubsequences(str, cur string, n int) {
	if N+1 == n {
		if len(cur) >= minLen {
			words[cur] = true
		}
	} else {
		generateSubsequences(str, cur, n+1)
		generateSubsequences(str, cur+string(str[n]), n+1)
	}
}

func wordsExists(file io.Reader) {
	scanner := bufio.NewScanner(file)

	validWords := make([]string, 0)

	for scanner.Scan() {
		validWords = append(validWords, scanner.Text())
	}

	for word := range words {
		var wordFound bool
		for _, w := range validWords {
			if w == word {
				wordFound = true
			}
		}

		if !wordFound {
			words[word] = false
		}
	}
}

func printWords() {
	sortedWords := make([][]string, N+1)

	for w, _ := range words {
		if words[w] {
			n := len(w)
			sortedWords[n] = append(sortedWords[n], w)
		}
	}
	fmt.Print("\r")
	for _, sw := range sortedWords {
		for _, w := range sw {
			fmt.Println(w)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: word <charset> <size>")
		log.Fatalln("No charset provided")
	} else if len(os.Args) < 3 {
		fmt.Println("Usage: word <charset> <size>")
		log.Fatalln("No min size provided")
	}
	charset = os.Args[1]
	N = len(charset) - 1

	if len(os.Args) > 1 {
		minLen, _ = strconv.Atoi(os.Args[2])
		if minLen > N {
			log.Fatalln("Min lenght greater than charset size")
		}
	}

	fmt.Printf("Minimum word length: %d\n", minLen)
	generatePermutations(strings.Split(charset, ""), 0)
	fmt.Printf("Found %d permutations\n", len(permutations))
	// fmt.Println(permutations)

	for p, _ := range permutations {
		generateSubsequences(p, "", 0)
	}

	N = len(words)
	fmt.Printf("Total sequences to check: %d\n", N)

	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	wordsExists(file)
	fmt.Println()

	printWords()
}
