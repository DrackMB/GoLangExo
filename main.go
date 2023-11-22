package main

import (
	"encoding/json"
	"estiam/dictionary"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	fmt.Println("Starting dictionary server on port 8080")

	d, err := dictionary.New("dictionary.json")
	if err != nil {
		fmt.Println("Error creating dictionary:", err)
		return
	}
	done := make(chan error)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		go func() {
			actionAdd(d, word, definition, done)
		}()

		if err := <-done; err != nil {
			fmt.Println("Error adding entry:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Entry added successfully"))
	}).Methods("POST")

	r.HandleFunc("/define/{word}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]

		result, err := actionDefine(d, word, done)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		jsonBytes, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonBytes)

	}).Methods("GET")

	r.HandleFunc("/remove/{word}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]
		go func() {
			actionRemove(d, word, done)
		}()
		if err := <-done; err != nil {
			fmt.Println("Error removing entry:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Word removed successfully."))
	}).Methods("DELETE")

	r.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		_, entrie := actionList(d)
		jsonBytes, err := json.Marshal(entrie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonBytes)
	}).Methods("GET")

	http.ListenAndServe(":8083", r)
}
func actionAdd(d *dictionary.Dictionary, word string, definition string, done chan<- error) {
	d.Add(word, definition, done)
}

func actionDefine(d *dictionary.Dictionary, word string, done chan<- error) (dictionary.Entry, error) {
	res, err := d.Get(word)
	return res, err
}

func actionRemove(d *dictionary.Dictionary, word string, done chan<- error) {
	d.Remove(word, done)
	fmt.Println("Word removed successfully.")
}

func actionList(d *dictionary.Dictionary) ([]string, map[string]dictionary.Entry) {
	words, entries := d.List()
	return words, entries
}
