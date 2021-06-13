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

	fmt.Println("Start", issueName)
	copier.Copy(payload)
	issueName := fmt.Sprintf("Issue{Publisher %d, Publication %d, Issue %d}", payload.Publisher, payload.Publication, payload.Issue)
	fmt.Fprintf(w, "OK")
	fmt.Println("Done", issueName)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", copyIssue).Methods("POST")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func main() {
	handleRequests()
}
