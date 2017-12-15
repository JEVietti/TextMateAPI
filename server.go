package main

import (
	"log"
	"net/http"
	// Imports the Google Cloud Natural Language API client package.
	language "cloud.google.com/go/language/apiv1"
	"golang.org/x/net/context"
)

var ctx context.Context
var client *language.Client

func main() {
	ctx = context.Background()
	var err error
	// Creates a client.
	client, err = language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	router := NewRouter()

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", router))

}
