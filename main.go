package main

import (
	"bufio"
	"encoding/json"
	"estiam/dictionary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Use(logMiddleware)
	r.Use(authenticationMiddleware)
	fmt.Println("Starting dictionary server on port 8080")

	d, err := dictionary.New("dictionary.json")
	if err != nil {
		fmt.Println("Error creating dictionary:", err)
		return
	}
	done := make(chan error)
	reader := bufio.NewReader(os.Stdin)

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
			d.Add(word, definition, done)
		}()

		if err := <-done; err != nil {
			fmt.Println("Error adding entry:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Entry added successfully"))
	}).Methods("POST")

	r.HandleFunc("/define", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("Enter word to define: ")
		word, _ := reader.ReadString('\n')
		word = strings.TrimSpace(word)

		entry, err := d.Get(word)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Def: %s\n", entry)
	})
	r.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			fmt.Print("Enter word to remove: ")
			word, _ := reader.ReadString('\n')
			word = strings.TrimSpace(word)

			d.Remove(word, done)
			fmt.Println("Word removed successfully.")
		}()
		if err := <-done; err != nil {
			fmt.Println("Error removing entry:", err)
		}
	})
	r.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		words, entries := d.List()
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
	})

	http.ListenAndServe(":8083", r)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Open the log file
		logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()

		// Set the log output to the file
		log.SetOutput(logFile)

		// Write log entry
		log.Println("Incoming request from", r.RemoteAddr, r.Method, r.URL.Path)

		// Forward request to the next handler
		next.ServeHTTP(w, r)
	})

}

var validTokens = map[string]bool{
	"token1": true,
	"token2": true,
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the request header
		token := r.Header.Get("Authorization")

		// Validate the token
		if !isValidToken(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Forward request to the next handler
		next.ServeHTTP(w, r)
	})
}

func isValidToken(token string) bool {
	// Check if the token is valid
	if _, ok := validTokens[token]; !ok {
		return false
	}

	return true
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write log entry
		log.Println("Incoming request from", r.RemoteAddr, r.Method, r.URL.Path)

		// Validate the token
		if !isValidToken(r.Header.Get("Authorization")) {
			log.Println("Invalid token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Forward request to the next handler
		next.ServeHTTP(w, r)
	})
}
