package main

import (
	"context"
	"encoding/json"
	"estiam/dictionary"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func AddWord(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Automatically cancel the context on function exit

	done := make(chan error)

	go func() {
		defer close(done) // Close the channel when the operation finishes

		body, err := io.ReadAll(r.Body)
		if err != nil {
			// Handle error reading the request body
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var entry dictionary.Entry
		err = json.Unmarshal(body, &entry)
		if err != nil {
			// Handle error unmarshaling JSON data
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Add the word to the dictionary asynchronously
		err = d.Add(ctx, entry, done)
		if err != nil {
			done <- err
			return
		}

		// Signal completion
		done <- nil
	}()

	// Wait for the goroutine to signal completion before sending the response
	select {
	case err := <-done:
		if err != nil {
			// Handle error adding word to dictionary
			fmt.Println("Error adding entry:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send success response
		w.Write([]byte("Entry added successfully"))
		w.WriteHeader(http.StatusOK)

	case <-ctx.Done():
		// Handle context cancellation
		w.WriteHeader(http.StatusRequestTimeout)
	}
}

func DeleteWord(w http.ResponseWriter, r *http.Request) {
	done := make(chan error)

	word := mux.Vars(r)["word"]
	_, err := d.Get(word)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	go func() {
		d.Remove(word, done)
	}()
	if err := <-done; err != nil {
		fmt.Println("Error removing entry:", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Word removed successfully."))

}

func GetALL(w http.ResponseWriter, r *http.Request) {
	words, err := d.List()
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(words)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func GetDefin(w http.ResponseWriter, r *http.Request) {
	word := mux.Vars(r)["word"]

	entry, err := d.Get(word)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonData, err := json.Marshal(entry)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
