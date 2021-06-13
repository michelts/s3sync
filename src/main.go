package main

import (
	"copier"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func copyIssue(w http.ResponseWriter, r *http.Request) {
	var payload copier.Issue

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	issueName := fmt.Sprintf("Issue{Publisher %d, Publication %d, Issue %d}", payload.Publisher, payload.Publication, payload.Issue)
	fmt.Println("Start", issueName)
	copier.Copy(payload)
	fmt.Fprintf(w, "OK")
	fmt.Println("Done", issueName)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", copyIssue).Methods("POST")
	router.HandleFunc("/ht/", healthCheck).Methods("GET")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func main() {
	handleRequests()
}
