package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

type Users struct {
	Data string
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	var u Users

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Do something with the Person struct...
	fmt.Fprintf(w, "Person: %+v", u)
}
