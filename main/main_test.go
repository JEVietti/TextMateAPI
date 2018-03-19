package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	language "cloud.google.com/go/language/apiv1"
	"github.com/jevietti/textmate/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestAvgLength(t *testing.T) {
	t.Log("Testing Average Length of Message Calculation")
	expected := float32(23.0)
	actual := CalculateAvgLength(models.RatingInput{"This is a test message.", 1})
	if actual != expected {
		t.Errorf("Test failed: expected: %.1f vs actual: %.1f", expected, actual)
	}
}

func TestTodos(t *testing.T) {
	var todos models.Todos
	request, _ := http.NewRequest("GET", "/todos", nil)
	response := httptest.NewRecorder()
	NewRouter().ServeHTTP(response, request)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &todos)
	fmt.Println(todos)
	assert.Equal(t, false, todos[0].Completed, "First Todo not Complete")
	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestEmptyGetScore(t *testing.T) {
	var ratingResult models.Ratings
	ctx = context.Background()
	var err error
	// Creates a client.
	client, err = language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	data := url.Values{}
	data.Set("Input", "")
	data.Add("Count", "0")

	r, _ := http.NewRequest("POST", "/api/sentiment", strings.NewReader(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	NewRouter().ServeHTTP(response, r)
	//fmt.Println(response)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &ratingResult)
	fmt.Println(ratingResult.Length)
	fmt.Println(ratingResult.Sentiment.Score)

	for index, element := range ratingResult.Sentiment.SentimentStruct.GetSentences() {
		fmt.Printf("Sentence %x : %s has a sentiment score of %.1f\n", index, element.GetText().GetContent(), element.GetSentiment().GetScore())
	}
	assert.Equal(t, 100, ratingResult.RatingScore, "Empty Message starts at 100 score.")
}

func TestNeutralGetScore(t *testing.T) {
	var ratingResult models.Ratings
	ctx = context.Background()
	var err error
	// Creates a client.
	client, err = language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	data := url.Values{}
	data.Set("Input", "")
	data.Add("Count", "2")

	r, _ := http.NewRequest("POST", "/api/sentiment", strings.NewReader(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	NewRouter().ServeHTTP(response, r)
	//fmt.Println(response)
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &ratingResult)
	fmt.Println(ratingResult.Length)
	fmt.Println(ratingResult.Sentiment.Score)

	for index, element := range ratingResult.Sentiment.SentimentStruct.GetSentences() {
		fmt.Printf("Sentence %x : %s has a sentiment score of %.1f\n", index, element.GetText().GetContent(), element.GetSentiment().GetScore())
	}
	assert.Equal(t, 100, ratingResult.RatingScore, "Empty Message starts at 100 score.")
}
