package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/jevietti/textmate/models"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

var decoder = schema.NewDecoder()

func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	todos := models.Todos{
		models.Todo{Name: "Write presentation"},
		models.Todo{Name: "Host meetup"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

//Get Score is the handler for calculating the relationship between two users using many
//child functions such as GetRatings for their sentiment scores,
func GetScore(w http.ResponseWriter, r *http.Request) {
	var score models.Ratings

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("r.PostForm", r.PostForm)
	log.Println("r.Form", r.Form)

	var input models.RatingInput
	//fmt.Println(r.PostForm)
	// r.PostForm is a map of our POST form values
	err = decoder.Decode(&input, r.PostForm)

	if err != nil {
		// Handle error
		panic(err)
	}

	avgLengthChannel := make(chan float32)
	sentimentChannel := make(chan *languagepb.AnalyzeSentimentResponse)

	go GetAverageLength(input, avgLengthChannel)
	go GetSentiment(input, sentimentChannel)

	avgLength, _ := <-avgLengthChannel
	sentiment, _ := <-sentimentChannel

	for index, element := range sentiment.GetSentences() {
		fmt.Printf("Sentence %x : %s has a sentiment score of %.1f\n", index, element.GetText().GetContent(), element.GetSentiment().GetScore())
	}

	if sentiment.DocumentSentiment.Score >= 0 {
		fmt.Printf("Sentiment: positive  %.1f.\n", sentiment.DocumentSentiment.Score)
	} else {
		fmt.Printf("Sentiment: negative %.1f.\n", sentiment.DocumentSentiment.Score)
	}

	score.Length = input.Count
	score.RatingScore = CalculateRelationshipScore(avgLength, sentiment.DocumentSentiment.GetScore(), sentiment.DocumentSentiment.GetMagnitude())
	score.Sentiment = models.SentimentScore{sentiment, sentiment.DocumentSentiment.GetScore(), "nil", "nil"}

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

func CalculateAvgLength(input models.RatingInput) float32 {
	text := input.Input
	fmt.Println(text)
	return float32(len(text)) / float32(input.Count)
}

//
func CalculateRelationshipScore(avgLength float32, sentimentScore float32, sentimentMagnitude float32) float32 {
	fmt.Printf("Calculating Relationship score based on average length of %.1f, sentiment score of %.1f, and magnitude of %.1f.\n", avgLength, sentimentScore, sentimentMagnitude)
	return 100 / ((sentimentScore / sentimentMagnitude) * avgLength)
}

// GetSentiment gets the Sentiment Rating Score for a list of text or single messages
// this is done through sentiment analysis prodvided by Google Cloud NLP Sentiment Analysis
func GetSentiment(input models.RatingInput, sentimentScore chan *languagepb.AnalyzeSentimentResponse) {

	// Sets the text to analyze.
	text := input.Input

	// Detects the sentiment of the text.
	sentiment, err := client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
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
