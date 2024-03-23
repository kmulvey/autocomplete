package dictionary

import (
	"encoding/csv"
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/tj/assert"
)

func TestInsert(t *testing.T) {
	t.Parallel()

	var dictionary = NewDictionary()
	dictionary.InsertWord("grape")
	dictionary.InsertWord("maple")
	dictionary.InsertWord("mince")
	dictionary.InsertWord("mini")
	dictionary.InsertWord("miniature")
	assert.True(t, dictionary.Search("grape"))
	assert.True(t, dictionary.Search("maple"))
	assert.True(t, dictionary.Search("mince"))
	assert.True(t, dictionary.Search("mini"))
	assert.True(t, dictionary.Search("miniature"))
	assert.False(t, dictionary.Search("acrobat"))
	assert.False(t, dictionary.Search("destiny"))
	assert.False(t, dictionary.Search("merger"))
	assert.NotNil(t, dictionary.trieNode.nodeFromPrefix("grape"))
	assert.Nil(t, dictionary.trieNode.nodeFromPrefix("acrobat"))

	assert.NoError(t, dictionary.PopulateFromCSV("./english.csv"))
}

func TestPopulateFromCSV(t *testing.T) {
	t.Parallel()

	var dictionary = NewDictionary()
	assert.NoError(t, dictionary.PopulateFromCSV("./english.csv"))
	assert.True(t, dictionary.Search("grape"))
	assert.True(t, dictionary.Search("maple"))
	assert.True(t, dictionary.Search("pear"))
	assert.False(t, dictionary.Search("mom"))
	assert.False(t, dictionary.Search("mini"))
	assert.False(t, dictionary.Search("jennifer"))
	assert.NotNil(t, dictionary.trieNode.nodeFromPrefix("grape"))
	assert.Nil(t, dictionary.trieNode.nodeFromPrefix("jennifer"))
}

func TestCollect(t *testing.T) {
	t.Parallel()

	var dictionary = NewDictionary()
	dictionary.InsertWord("grape")
	dictionary.InsertWord("maple")
	dictionary.InsertWord("mince")
	dictionary.InsertWord("mini")
	dictionary.InsertWord("miniature")

	var words = dictionary.Collect("")
	assert.EqualValues(t, []string{"grape", "maple", "mince", "mini", "miniature"}, words)
}

// TestEnglishCSV tests that the whole dictionary can be added without error
func TestEnglishCSV(t *testing.T) {
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

	var dictionary = NewDictionary()
	assert.NoError(t, dictionary.PopulateFromCSV("./english.csv"))
	var words = dictionary.Collect("")
	sort.Strings(words)
	assert.Equal(t, 111723, len(words))

	for i, word := range uniqWordsArr {
		if word != words[i] {
			assert.Fail(t, "they should all match")
		}
	}
}
