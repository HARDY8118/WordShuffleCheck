package main

import (
	"fmt"
	"log"
	"net/http"
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

const API_PATH = "https://api.dictionaryapi.dev/api/v2/entries/en/%s"

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

func wordExists(word string) bool {
	resp, err := http.Get(fmt.Sprintf(API_PATH, word))

	if err != nil {
		fmt.Println("Failed to connect")
		log.Fatalln(err)
	}

	wordsChecked += 1
	fmt.Printf("\r%d/%d [%0.4f%%]", wordsChecked, N, 100*(float64(wordsChecked)/float64(N)))

	switch resp.StatusCode {
	case 200:
		return true
	case 404:
		// delete(words, word)
		words[word] = false
		return false
	case 429:
		log.Fatalln("Too many requests made")
		return false
	default:
		log.Panicf("API Call failed for word '%s' with status code: %d", word, resp.StatusCode)
		return false
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
		log.Fatalln("No charset provided")
	}
	charset = os.Args[1]
	N = len(charset) - 1

	if len(os.Args) > 1 {
		minLen, _ = strconv.Atoi(os.Args[2])
	}

	fmt.Printf("Minimum word length: %d\n", minLen)
	generatePermutations(strings.Split(charset, ""), 0)
	fmt.Printf("Found %d permutations\n", len(permutations))

	for p, _ := range permutations {
		generateSubsequences(p, "", 0)
	}

	N = len(words)
	fmt.Printf("Total sequences to check: %d\n", N)

	for w, _ := range words {
		wordExists(w)
	}
	fmt.Println()

	printWords()
}
