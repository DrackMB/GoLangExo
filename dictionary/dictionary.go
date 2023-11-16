package dictionary

import "fmt"

type Entry struct {
	Def string
}

func (e Entry) String() string {

	return e.Def
}

type Dictionary struct {
	entries map[string]Entry
}

func New() *Dictionary {

	return &Dictionary{entries: make(map[string]Entry)}
}

func (d *Dictionary) Add(word string, definition string) {
	entry := Entry{Def: definition}
	d.entries[word] = entry
}

func (d *Dictionary) Get(word string) (Entry, error) {

	entry, exists := d.entries[word]
	if !exists {
		return Entry{}, fmt.Errorf("Word not found: %s", word)
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) {
	delete(d.entries, word)
}

func (d *Dictionary) List() ([]string, map[string]Entry) {

	words := make([]string, 0, len(d.entries))
	for word := range d.entries {
		words = append(words, word)
	}
	return words, d.entries
}
