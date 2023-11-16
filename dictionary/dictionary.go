package dictionary

import (
	"encoding/json"
	"fmt"
	"os"
)

type Entry struct {
	Definition string
}

type Dictionary struct {
	filePath string
	entries  map[string]Entry
}

func New(filePath string) (*Dictionary, error) {
	d := &Dictionary{filePath: filePath, entries: make(map[string]Entry)}
	if err := d.loadFromFile(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Dictionary) Add(word string, definition string) error {
	d.entries[word] = Entry{Definition: definition}
	return d.saveToFile()
}

func (d *Dictionary) Get(word string) (Entry, error) {
	entry, exists := d.entries[word]
	if !exists {
		return Entry{}, fmt.Errorf("Word not found: %s", word)
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) error {
	delete(d.entries, word)
	return d.saveToFile()
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	words := make([]string, 0, len(d.entries))
	for word := range d.entries {
		words = append(words, word)
	}
	return words, d.entries
}

func (d *Dictionary) loadFromFile() error {
	fileData, err := os.ReadFile(d.filePath)
	if err != nil {
		return nil
	}

	if len(fileData) == 0 {
		d.entries = make(map[string]Entry)
		return nil
	}

	if err := json.Unmarshal(fileData, &d.entries); err != nil {
		return fmt.Errorf("Error JSON: %v", err)
	}

	return nil
}

func (d *Dictionary) saveToFile() error {
	fileData, err := json.MarshalIndent(d.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("Error JSON: %v", err)
	}

	if err := os.WriteFile(d.filePath, fileData, 0644); err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}
