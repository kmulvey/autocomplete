package main

import (
	"encoding/csv"
	"io"
	"os"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/tj/assert"
)

func TestInsert(t *testing.T) {
	t.Parallel()

	var trie = NewTrie()
	trie.InsertWord("mom")
	trie.InsertWord("mince")
	trie.InsertWord("mini")
	trie.InsertWord("moose")
	trie.InsertWord("maple")
	trie.InsertWord("grape")
	trie.InsertWord("pear")
	assert.NotNil(t, trie.Search("mom"))
	assert.NotNil(t, trie.Search("pear"))

	trie.InsertDictionaryFromCSV("./english.csv")
	assert.NotNil(t, trie.Search("grape"))
	assert.NotNil(t, trie.Search("maple"))
	assert.NotNil(t, trie.Search("mince"))
	assert.NotNil(t, trie.Search("moose"))
	assert.NotNil(t, trie.Search("pear"))

	var words = Collect(trie, "")
	sort.Strings(words)
	assert.True(t, slices.Contains(words, "grape"))
	assert.True(t, slices.Contains(words, "maple"))
	assert.True(t, slices.Contains(words, "mince"))
	assert.True(t, slices.Contains(words, "moose"))
	assert.True(t, slices.Contains(words, "pear"))

	var results = Autocomplete(trie, "apple")
	sort.Strings(results)
	assert.EqualValues(t, []string{"apple pie", "apple-faced", "apple-jack", "apple-john", "apple-squire"}, results)
}

func TestCSV(t *testing.T) {
	f, err := os.Open("./english.csv")
	assert.NoError(t, err)

	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	var header = true
	var uniqWordsMap = make(map[string]struct{})
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)

		if !header && isASCII(rec[0]) {
			uniqWordsMap[strings.ToLower(rec[0])] = struct{}{}
		}
		header = false
	}

	var uniqWordsArr = make([]string, len(uniqWordsMap))
	var i int
	for word := range uniqWordsMap {
		uniqWordsArr[i] = word
		i++
	}
	sort.Strings(uniqWordsArr)
	assert.Equal(t, 111723, len(uniqWordsArr))

	var trie = NewTrie()
	trie.InsertDictionaryFromCSV("./english.csv")
	var words = Collect(trie, "")
	sort.Strings(words)
	assert.Equal(t, 111723, len(words))

	for i, word := range uniqWordsArr {
		if word != words[i] {
			assert.Fail(t, "they should all match")
		}
	}
}
