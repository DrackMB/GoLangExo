package dictionary_test

import (
	"testing"

	"estiam/dictionary"

	"github.com/stretchr/testify/assert"
)

func TestAddInvalidWord(t *testing.T) {
	// Create a new dictionary
	dict := dictionary.New("test.json")

	// Try adding an invalid word (empty word)
	done := make(chan error)
	go dict.Add("", "Dif", done)
	err := <-done
	defer close(done)

	assert.Error(t, err)
	assert.EqualError(t, err, "Definition cannot be empty")

	// Try adding an invalid word (null word)
	err = dict.Add("Word", "", done)
	defer close(done)

	assert.Error(t, err)
	assert.EqualError(t, err, "Definition cannot be null")
}

func TestGetInvalidWord(t *testing.T) {
	// Create a new dictionary
	dict := dictionary.New("test.json")

	// Try getting an invalid word
	_, err := dict.Get("banana")
	assert.Error(t, err)
	assert.EqualError(t, err, "Word not found: banana")
}

func TestRemoveInvalidWord(t *testing.T) {
	// Create a new dictionary
	dict := dictionary.New("test.json")
	done := make(chan error)
	// Try removing an invalid word
	err := dict.Remove("banana", done)
	assert.Error(t, err)
	assert.EqualError(t, err, "Word not found: banana")
}

func TestListEmptyDictionary(t *testing.T) {
	// Create a new dictionary
	dict := dictionary.New("test.json")
	// List the words
	words, _, err := dict.List()
	assert.Nil(t, err)
	assert.Equal(t, words, []string{})
}
