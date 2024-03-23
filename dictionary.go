package dictionary

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
	"unicode"
)

type Dictionary struct {
	*trieNode
}

type trieNode struct {
	Letter    string
	Children  map[string]*trieNode
	EndOfWord bool
}

func (l trieNode) String() string {
	return l.Letter
}

func NewDictionary() *Dictionary {
	return &Dictionary{&trieNode{Children: make(map[string]*trieNode, 26)}}
}

func (d *Dictionary) InsertDictionaryFromCSV(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	csvReader := csv.NewReader(f)
	var header = true
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return f.Close()
		}

		// we dont want the file header or words with wierd characters
		if !header && isASCII(rec[0]) {
			d.InsertWord(rec[0])
		}
		header = false
	}

	return f.Close()
}

// InsertWord adds the word to the dictionary
func (d *Dictionary) InsertWord(word string) {

	word = strings.ToLower(word)

	var lastNode, found = d.trieNode.Children[string(word[0])]
	if !found {
		d.trieNode.Children[string(word[0])] = &trieNode{Letter: string(word[0]), Children: make(map[string]*trieNode, 2)}
		lastNode = d.trieNode.Children[string(word[0])]
	}

	if len(word) == 1 {
		lastNode.EndOfWord = true
		return
	}

	for i := 1; i < len(word); i++ {
		var currLetter = string(word[i])
		var endOfWord = i == len(word)-1

		if lastNode.Children == nil {
			lastNode.Children = make(map[string]*trieNode, 2)
		}

		if _, exists := lastNode.Children[currLetter]; !exists {
			lastNode.Children[currLetter] = &trieNode{Letter: currLetter, Children: make(map[string]*trieNode, 2), EndOfWord: endOfWord}
		}

		if !lastNode.Children[currLetter].EndOfWord { // if its already true dont reset it to false for a longer word
			lastNode.Children[currLetter].EndOfWord = endOfWord
		}
		lastNode = lastNode.Children[currLetter]
	}
}

// Search returns true if the given word is found in the trie
func (t *trieNode) Search(word string) bool {

	word = strings.ToLower(word)
	var lastNode = t.Children[string(word[0])]
	if lastNode == nil {
		return false
	}

	for i := 1; i < len(word); i++ {
		var currLetter = string(word[i])

		if _, found := lastNode.Children[currLetter]; !found {
			return false
		}

		lastNode = lastNode.Children[currLetter]
	}

	return lastNode.EndOfWord
}

// NodeFromPrefix is similar to Search() but returns a node
func (t *trieNode) NodeFromPrefix(word string) *trieNode {

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

	return lastNode
}

// Collect recursively gathers all words in the Trie
func Collect(node *trieNode, prefix string) []string {

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

func (d *Dictionary) Autocomplete(prefix string) []string {
	return autocomplete(d.trieNode, prefix)
}

func autocomplete(node *trieNode, prefix string) []string {
	var results = Collect(node.NodeFromPrefix(prefix), "")
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
