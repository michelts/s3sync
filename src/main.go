package main

import (
	"copier"
	"fmt"
	"log"
	"net/http"
)

func copyIssue(w http.ResponseWriter, r *http.Request) {
	copier.Copy()
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", copyIssue).Methods("POST")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func main() {
	handleRequests()
}
