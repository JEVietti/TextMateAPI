package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"

	language "cloud.google.com/go/language/apiv1"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	models "github.com/jevietti/textmate/models"
	"google.golang.org/appengine"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

var decoder = schema.NewDecoder()

func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	//fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	todos := Todos{
		Todo{Name: "Write presentation"},
		Todo{Name: "Host meetup"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

//
func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoID"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

//Get Score is the handler for calculating the relationship between two users using many
//child functions such as GetRatings for their sentiment scores,
func GetScore(w http.ResponseWriter, req *http.Request) {
	ctx = appengine.NewContext(req)
	var err error
	// Creates a client.
	client, err = language.NewClient(ctx)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	var score models.Ratings
	var input models.RatingInput

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&input)
	if err != nil {
		panic(err)
	}

	avgLengthChannel := make(chan float32, 1)
	sentimentChannel := make(chan *languagepb.AnalyzeSentimentResponse, 1)

	go GetAverageLength(input, avgLengthChannel)
	go GetSentiment(input, sentimentChannel)

	sentiment, _ := <-sentimentChannel
	avgLength, _ := <-avgLengthChannel

	score.Length = input.Count
	score.RatingScore = CalculateRelationshipScore(avgLength, sentiment.DocumentSentiment.GetScore(), sentiment.DocumentSentiment.GetMagnitude())
	score.Sentiment = models.SentimentScore{SentimentStruct: sentiment, Score: sentiment.DocumentSentiment.GetScore(), BestMessage: "nil", WorstMessage: "nil"}
	////fmt.Printf("Relationship score is %.1f.\n", score.RatingScore)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(score); err != nil {
		panic(err)
	}

}

// GetAverageLength is a helper function for the RelationshipAlgorithm
// to find the avg length of a message for the set of messages sent
func GetAverageLength(input models.RatingInput, avgLength chan float32) {

	avgLength <- CalculateAvgLength(input)
	close(avgLength)
}

//
func CalculateAvgLength(input models.RatingInput) float32 {
	text := strings.Join(input.Input, " ")
	return float32(len(text)) / float32(input.Count)
}

//
func CalculateRelationshipScore(avgLength float32, sentimentScore float32, sentimentMagnitude float32) float32 {
	////fmt.Printf("Calculating Relationship score based on average length of %.1f, sentiment score of %.1f, and magnitude of %.1f.\n", avgLength, sentimentScore, sentimentMagnitude)
	score := (1 - (avgLength / 100)) * 25
	sentiment := (sentimentScore * sentimentMagnitude) + sentimentMagnitude
	if sentiment < -0.5 {
		score += sentiment
	} else {
		score += float32(math.Abs(float64(sentiment)) * 75)
	}
	if score >= 100.0 {
		return 100.0
	} else if score <= 0.0 {
		return 0.0
	}
	return score
}

// GetSentiment gets the Sentiment Rating Score for a list of text or single messages
// this is done through sentiment analysis prodvided by Google Cloud NLP Sentiment Analysis
func GetSentiment(input models.RatingInput, sentimentScore chan *languagepb.AnalyzeSentimentResponse) {
	// Sets the text to analyze.
	text := input.Input
	messages := strings.Join(text, " ")
	// Detects the sentiment of the text.
	sentiment, err := client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: messages,
			},
			Type:     languagepb.Document_PLAIN_TEXT,
			Language: "en",
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}
	sentimentScore <- sentiment
	close(sentimentScore)
}
