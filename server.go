package main

import (

	// Imports the Google Cloud Natural Language API client package.

	"log"
	"net/http"

	language "cloud.google.com/go/language/apiv1"
	"golang.org/x/net/context"
)

var ctx context.Context
var client *language.Client

//Google Server
/*
func init() {
	router := NewRouter()

	http.Handle("/", router)
}
*/
//Normal Server

func main() {

	router := NewRouter()

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", router))

}
