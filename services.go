package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func AddWord(w http.ResponseWriter, r *http.Request) {
	done := make(chan error)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(strings.NewReader(string(body)))
	var data map[string]string
	err = decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	word, ok := data["word"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing word field in request body"))
		return
	}

	definition, ok := data["definition"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing definition field in request body"))
		return
	}
	go d.Add(word, definition, done)

	if err := <-done; err != nil {
		fmt.Println("Error adding entry:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Entry added successfully"))
}

func DeleteWord(w http.ResponseWriter, r *http.Request) {
	done := make(chan error)

	word := mux.Vars(r)["word"]
	fmt.Println(word)
	go d.Remove(word, done)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Word removed successfully."))

	if err := <-done; err != nil {
		fmt.Println("Error removing entry:", err)
	}
}

func GetALL(w http.ResponseWriter, r *http.Request) {
	words, entries, _ := d.List()
	fmt.Println(words)
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(entries)
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
		fmt.Println(err)
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
