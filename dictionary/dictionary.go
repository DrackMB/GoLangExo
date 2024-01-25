package dictionary

import (
	"context"
	"estiam/db"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Entry struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type Dictionary struct {
	Db *redis.Client
}

func New() *Dictionary {
	return &Dictionary{
		Db: db.DatabaseConnect(),
	}
}

func (d *Dictionary) Add(ctx context.Context, entry Entry, done chan<- error) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure the context is canceled to avoid leaks
	fmt.Printf(entry.Definition)
	// Validate the word and definition
	if entry.Word == "" || entry.Definition == "" {
		err := fmt.Errorf("Word or definition cannot be empty")
		return err
	}

	// Create a goroutine to handle the operation and signal completion
	go func() {
		defer close(done) // Close the channel when the operation finishes
		fmt.Println(entry)
		// Add the word to the dictionary
		err := d.Db.Set(ctx, entry.Word, entry.Definition, 0).Err()
		if err != nil {
			done <- fmt.Errorf("Error adding word to database: %s", err)
			return
		}

		done <- nil // Signal completion
	}()

	// Wait for the goroutine to signal completion before returning
	select {
	case <-ctx.Done():
		return fmt.Errorf("Operation timed out")
	}
}

func (d *Dictionary) Get(word string) (Entry, error) {
	val, err := d.Db.Get(context.Background(), word).Result()
	if err == redis.Nil {
		return Entry{}, fmt.Errorf("Word not found: %s", word)
	} else if err != nil {
		log.Fatal(err)
	}
	return Entry{Word: word, Definition: val}, nil
}

func (d *Dictionary) Remove(word string, done chan<- error) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel() // Automatically cancel the context on function exit

	// Wait for the context to be canceled or the operation to complete
	go func() {
		err := d.Db.Del(ctx, word).Err()
		if err != nil && err != context.DeadlineExceeded {
			done <- fmt.Errorf("error deleting word: %s", word)
		} else {
			done <- nil
		}
	}()

	// Wait for the goroutine to signal completion
	select {
	case <-ctx.Done():
		return fmt.Errorf("Word not found: %s", word)
	}
}

func (d *Dictionary) List() ([]Entry, error) {
	ctx := context.Background()

	keys, err := d.Db.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	var entries []Entry
	for _, key := range keys {
		fmt.Printf(key)
		value, err := d.Db.Get(ctx, key).Result()
		fmt.Printf(value)
		if err != nil {
			return nil, err
		}

		var entry Entry
		entry.Definition = value
		entry.Word = key
		entries = append(entries, entry)
	}

	return entries, nil
}
