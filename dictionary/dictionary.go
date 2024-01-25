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

func New(filePath string) *Dictionary {
	d := &Dictionary{filePath: filePath, entries: make(map[string]Entry)}
	if err := d.loadFromFile(); err != nil {
		return nil
	}
	return d
}

func (d *Dictionary) Add(word string, definition string, done chan<- error) error {
	// Validate the word and definition
	if word == "" || definition == "" {
		err := fmt.Errorf("Word or definition cannot be empty")
		done <- err
		return err
	}

	// Add the word to the dictionary
	d.entries[word] = Entry{Definition: definition}

	err := d.saveToFile()
	done <- err
	return nil
}

func (d *Dictionary) Get(word string) (Entry, error) {
	entry, exists := d.entries[word]
	if !exists {
		return Entry{}, fmt.Errorf("Word not found: %s", word)
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string, done chan<- error) error {
	entry, _ := d.Get(word)
	if entry != (Entry{}) {
		delete(d.entries, word)
		err := d.saveToFile()
		done <- err
		return nil
	}
	err := fmt.Errorf("Word not found: %s", word)
	return err
}

func (d *Dictionary) List() ([]string, map[string]Entry, error) {
	words := make([]string, 0, len(d.entries))
	for word := range d.entries {
		words = append(words, word)
	}
	return words, d.entries, nil
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
		return fmt.Errorf("Error unmarshalling JSON: %v", err)
	}

	return nil
}

func (d *Dictionary) saveToFile() error {
	fileData, err := json.MarshalIndent(d.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshalling JSON: %v", err)
	}

	if err := os.WriteFile(d.filePath, fileData, 0644); err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}
