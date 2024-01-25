package dictionary

import (
	"context"
	"testing"

	"estiam/db"

	"github.com/stretchr/testify/assert"
)

func TestAddInvalidWord(t *testing.T) {
	// Create a new dictionary
	dict := &Dictionary{
		Db: db.DatabaseConnect(),
	}
	// Try adding an invalid word (empty word)
	done := make(chan error)
	entry := Entry{Word: "", Definition: "test"}
	ctx, _ := context.WithCancel(context.Background())
	err := dict.Add(ctx, entry, done)

	assert.Error(t, err)
	assert.EqualError(t, err, "Word or definition cannot be empty")
}

func TestGetInvalidWord(t *testing.T) {
	// Create a new dictionary
	dict := &Dictionary{
		Db: db.DatabaseConnect(),
	}

	// Try getting an invalid word
	_, err := dict.Get("banana")
	assert.Error(t, err)
	assert.EqualError(t, err, "Word not found: banana")
}

func TestRemoveInvalidWord(t *testing.T) {
	// Create a new dictionary
	dict := &Dictionary{
		Db: db.DatabaseConnect(),
	}
	done := make(chan error)
	// Try removing an invalid word
	err := dict.Remove("banana", done)
	assert.Error(t, err)
	assert.EqualError(t, err, "Word not found: banana")
}

func TestListEmptyDictionary(t *testing.T) {
	// Create a new dictionary
	dict := &Dictionary{
		Db: db.DatabaseConnect(),
	}
	// List the words
	words, err := dict.List()
	assert.Nil(t, err)
	assert.Equal(t, words, []Entry{{Word: "mb23", Definition: "mb345"}})
}
