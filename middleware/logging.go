package middleware

import (
	"log"
	"net/http"
	"os"
)

func LogMiddleware(next http.Handler) http.Handler {
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
