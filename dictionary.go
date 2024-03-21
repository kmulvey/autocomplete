package autocomplete

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
	"unicode"
)

type Trie struct {
	Letter    string
	Children  map[string]*Trie
	EndOfWord bool
}

func (l Trie) String() string {
	return l.Letter
}

func NewTrie() *Trie {
	var trie = Trie{Children: make(map[string]*Trie, 26)}

	return &trie
}

func (t *Trie) InsertDictionaryFromCSV(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	var header = true
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// do something with read line
		if !header && isASCII(rec[0]) {
			t.InsertWord(rec[0])
		}
		header = false
	}

	return nil
}

func (t *Trie) InsertWord(word string) {

	word = strings.ToLower(word)

	var lastNode, found = t.Children[string(word[0])]
	if !found {
		t.Children[string(word[0])] = &Trie{Letter: string(word[0]), Children: make(map[string]*Trie, 2)}
		lastNode = t.Children[string(word[0])]
	}

	if len(word) == 1 {
		lastNode.EndOfWord = true
		return
	}

	for i := 1; i < len(word); i++ {
		var currLetter = string(word[i])
		var endOfWord = i == len(word)-1

		if lastNode.Children == nil {
			lastNode.Children = make(map[string]*Trie, 2)
		}

		if _, exists := lastNode.Children[currLetter]; !exists {
			lastNode.Children[currLetter] = &Trie{Letter: currLetter, Children: make(map[string]*Trie, 2), EndOfWord: endOfWord}
		}

		if !lastNode.Children[currLetter].EndOfWord { // if its already true dont reset it to false for a longer word
			lastNode.Children[currLetter].EndOfWord = endOfWord
		}
		lastNode = lastNode.Children[currLetter]
	}
}

func (t *Trie) Search(word string) *Trie {

	word = strings.ToLower(word)
	var lastNode = t.Children[string(word[0])]
	if lastNode == nil {
		return nil
	}

	for i := 1; i < len(word); i++ {
		var currLetter = string(word[i])

		if _, found := lastNode.Children[currLetter]; !found {
			return nil
		}

		lastNode = lastNode.Children[currLetter]
	}

	if lastNode.EndOfWord {
		return lastNode
	}

	return nil
}

// Collect recursively gathers all words in the Trie
func Collect(node *Trie, prefix string) []string {

	if node == nil {
		return nil
	}
	var words []string
	if node.EndOfWord {
		words = append(words, prefix)
	}

	for char, node := range node.Children {
		words = append(words, Collect(node, prefix+char)...)
	}

	return words
}

func Autocomplete(node *Trie, prefix string) []string {
	var results = Collect(node.Search(prefix), "")
	if len(results) == 0 {
		return results
	}

	if results[0] == "" {
		results = results[1:]
	}

	for i, word := range results {
		results[i] = prefix + word
	}

	return results
}

func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}
